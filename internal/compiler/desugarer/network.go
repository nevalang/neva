package desugarer

import (
	"errors"
	"fmt"
	"maps"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

type handleNetworkResult struct {
	desugaredConnections []src.Connection     // desugared network
	constsToInsert       map[string]src.Const // constants that needs to be inserted in to make desugared network work
	nodesToInsert        map[string]src.Node  // nodes that needs to be inserted in to make desugared network work
	nodesPortsUsed       nodePortsMap         // to find unused to create virtual del connections
}

func (d Desugarer) handleNetwork(
	net []src.Connection,
	nodes map[string]src.Node,
	scope src.Scope,
) (handleNetworkResult, *compiler.Error) {
	nodesToInsert := map[string]src.Node{}
	constsToInsert := map[string]src.Const{}
	nodesPortsUsed := newNodePortsMap()

	desugaredConnections, err := d.desugarConnections(
		net,
		nodesPortsUsed,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
	)
	if err != nil {
		return handleNetworkResult{}, err
	}

	return handleNetworkResult{
		desugaredConnections: desugaredConnections,
		nodesPortsUsed:       nodesPortsUsed,
		constsToInsert:       constsToInsert,
		nodesToInsert:        nodesToInsert,
	}, nil
}

func (d Desugarer) desugarConnections(
	net []src.Connection,
	nodePortsUsed nodePortsMap,
	scope src.Scope,
	nodes map[string]src.Node,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) ([]src.Connection, *compiler.Error) {
	desugaredConnections := make([]src.Connection, 0, len(net))

	for _, conn := range net {
		result, err := d.desugarConnection(
			conn,
			nodePortsUsed,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return nil, err
		}

		desugaredConnections = append(desugaredConnections, result.connectionToReplace)
		desugaredConnections = append(desugaredConnections, result.connectionsToInsert...)
	}

	return desugaredConnections, nil
}

type desugarConnectionResult struct {
	connectionToReplace src.Connection
	connectionsToInsert []src.Connection
}

// desugarConnection modifies given nodesToInsert, constsToInsert and usedNodePorts
// it also returns connection to replace the original one and other connections
// that were generated while desugared the original one.
func (d Desugarer) desugarConnection(
	conn src.Connection,
	nodePortsUsed nodePortsMap,
	scope src.Scope,
	nodes map[string]src.Node,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (desugarConnectionResult, *compiler.Error) {
	if conn.ArrayBypass != nil {
		nodePortsUsed.set(
			conn.ArrayBypass.SenderOutport.Node,
			conn.ArrayBypass.SenderOutport.Port,
		)
		nodePortsUsed.set(
			conn.ArrayBypass.ReceiverInport.Node,
			conn.ArrayBypass.ReceiverInport.Port,
		)
		return desugarConnectionResult{connectionToReplace: conn}, nil
	}

	return d.desugarNormalConnection(
		*conn.Normal,
		nodePortsUsed,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
	)
}

func (d Desugarer) desugarNormalConnection(
	normConn src.NormalConnection,
	nodePortsUsed nodePortsMap,
	scope src.Scope,
	nodes map[string]src.Node,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (desugarConnectionResult, *compiler.Error) {
	if len(normConn.SenderSide) > 1 {
		result, err := d.desugarMultipleSenders(
			normConn,
			nodesToInsert,
			constsToInsert,
			nodePortsUsed,
			scope,
			nodes,
		)
		if err != nil {
			return desugarConnectionResult{}, &compiler.Error{
				Err:      err,
				Location: &scope.Location,
				Meta:     &normConn.Meta,
			}
		}

		// original connection is replaced by multiple new ones
		return desugarConnectionResult{
			connectionsToInsert: result,
		}, nil
	}

	singleSender := normConn.SenderSide[0]

	// sender might be const, range or have selectors so need to desugar
	// new connections might be created so we need to insert them and replace original one
	desugarSenderResult, err := d.desugarSender(
		normConn,
		singleSender,
		scope,
		nodes,
		nodePortsUsed,
		nodesToInsert,
		constsToInsert,
	)
	if err != nil {
		return desugarConnectionResult{}, compiler.Error{
			Location: &scope.Location,
			Meta:     &singleSender.Meta,
		}.Wrap(err)
	}

	normConn = *desugarSenderResult.connToReplace.Normal
	connectionsToInsert := desugarSenderResult.connectionsToInsert

	// there will be deferred connections, no chains, just normal receivers
	desugaredReceivers := make([]src.ConnectionPortReceiver, 0, len(normConn.ReceiverSide.Receivers))

	// desugar unnamed receivers (if needed)
	for _, portReceiver := range normConn.ReceiverSide.Receivers {
		if portReceiver.PortAddr.Port == "" {
			var err *compiler.Error
			portReceiver, err = d.desugarUnnamedReceiver(portReceiver, scope, nodes)
			if err != nil {
				return desugarConnectionResult{}, err
			}
		}
		desugaredReceivers = append(desugaredReceivers, portReceiver)
		// we don't return yet because there might also be deferred connections
	}

	if normConn.ReceiverSide.DeferredConnections != nil {
		result, err := d.desugarDeferredConnections(
			normConn,
			nodes,
			scope,
		)
		if err != nil {
			return desugarConnectionResult{}, err
		}

		// desugaring of deferred connections is recursive process so its result must be merged with existing one
		nodePortsUsed.merge(result.nodesPortsUsed)
		maps.Copy(constsToInsert, result.constsToInsert)
		maps.Copy(nodesToInsert, result.nodesToInsert)
		// after desugaring of deferred connection we need to add new receivers and new connections
		desugaredReceivers = append(desugaredReceivers, result.receiversToInsert...)
		connectionsToInsert = append(connectionsToInsert, result.connsToInsert...)
	}

	// desugar fan-out if needed
	if len(desugaredReceivers) > 1 {
		result := d.desugarFanOut(desugaredReceivers, nodesToInsert)
		// replace all existing receivers with single one
		desugaredReceivers = []src.ConnectionPortReceiver{result.receiverToReplace}
		connectionsToInsert = append(connectionsToInsert, result.connectionsToInsert...)
	}

	if normConn.ReceiverSide.ChainedConnection != nil {
		chainedConn := normConn.ReceiverSide.ChainedConnection
		chainHead := chainedConn.Normal.SenderSide[0] // chain head is always single sender

		if chainHead.Const != nil {
			return desugarConnectionResult{}, &compiler.Error{
				Err:      errors.New("chained connection with const sender is not supported"),
				Location: &scope.Location,
				Meta:     &normConn.Meta,
			}
		}

		// before desugar chained connection itself, we need to figure receiver side
		// if chain head is range, then port will be "sig"
		// if it's unnamed port-addr, then it will be first inport found
		var receiverPortName string
		if chainHead.Range != nil {
			receiverPortName = "sig"
		} else if chainHead.PortAddr != nil {
			var firstInportName = chainHead.PortAddr.Port
			if chainHead.PortAddr.Port == "" {
				var err error
				firstInportName, err = getFirstInportName(scope, nodes, *chainHead.PortAddr)
				if err != nil {
					return desugarConnectionResult{}, &compiler.Error{Err: err}
				}
			}
			receiverPortName = firstInportName
		}

		// recursively desugar the chained connection
		desugarChainResult, err := d.desugarConnection(*chainedConn, nodePortsUsed, scope, nodes, nodesToInsert, constsToInsert)
		if err != nil {
			return desugarConnectionResult{}, err
		}

		desugaredHead := desugarChainResult.connectionToReplace.Normal.SenderSide[0] // chain head is always single sender

		// connect sender to chain head by adding receiver to current connection
		desugaredReceivers = append(desugaredReceivers, src.ConnectionPortReceiver{
			PortAddr: src.PortAddr{
				Node: desugaredHead.PortAddr.Node, // node from head (sender)
				Port: receiverPortName,            // but port that we found before desugaring
				Meta: chainHead.Meta,
			},
			Meta: chainedConn.Meta,
		})

		// we need to insert both conn to replace and to insert, example:
		// input = a -> b -> c -> d
		// sender = a
		// chain = b -> c -> d
		// output = { to replace = b -> c, to insert = c -> d }
		connectionsToInsert = append(connectionsToInsert, desugarChainResult.connectionToReplace)
		connectionsToInsert = append(connectionsToInsert, desugarChainResult.connectionsToInsert...)
	}

	return desugarConnectionResult{
		connectionToReplace: src.Connection{
			Normal: &src.NormalConnection{
				SenderSide: normConn.SenderSide,
				ReceiverSide: src.ConnectionReceiverSide{
					Receivers: desugaredReceivers,
				},
			},
			Meta: normConn.Meta,
		},
		connectionsToInsert: connectionsToInsert,
	}, nil
}

type desugarSenderResult struct {
	connToReplace       src.Connection
	connectionsToInsert []src.Connection
}

func (d Desugarer) desugarSender(
	normConn src.NormalConnection,
	sender src.ConnectionSender,
	scope src.Scope,
	nodes map[string]src.Node,
	nodePortsUsed nodePortsMap,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (desugarSenderResult, *compiler.Error) {
	// mark as used and handle unnamed port if needed
	if sender.PortAddr != nil {
		if sender.PortAddr.Port == "" {
			firstOutportName, err := getFirstOutportName(scope, nodes, *sender.PortAddr)
			if err != nil {
				return desugarSenderResult{}, &compiler.Error{Err: err}
			}

			sender = src.ConnectionSender{
				PortAddr: &src.PortAddr{
					Port: firstOutportName,
					Node: sender.PortAddr.Node,
					Idx:  sender.PortAddr.Idx,
					Meta: sender.PortAddr.Meta,
				},
				Selectors: sender.Selectors,
				Meta:      sender.Meta,
			}
		}

		nodePortsUsed.set(
			sender.PortAddr.Node,
			sender.PortAddr.Port,
		)
	}

	connectionsToInsert := []src.Connection{}

	// if conn has selectors, desugar them, replace original connection and insert what's needed
	if len(sender.Selectors) != 0 {
		desugarSelectorsResult, err := d.desugarStructSelectors(
			sender,
			normConn,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return desugarSenderResult{}, compiler.Error{
				Location: &scope.Location,
				Meta:     &sender.Meta,
			}.Wrap(err)
		}

		// generated connection might need desugaring itself
		connToInsertDesugarRes, err := d.desugarConnection(
			desugarSelectorsResult.connToInsert,
			nodePortsUsed,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return desugarSenderResult{}, err
		}

		connectionsToInsert = append(connectionsToInsert, connToInsertDesugarRes.connectionToReplace)
		connectionsToInsert = append(connectionsToInsert, connToInsertDesugarRes.connectionsToInsert...)

		// connection that replaces original one might need desugaring itself
		replacedConnDesugarRes, err := d.desugarConnection(
			desugarSelectorsResult.connToReplace,
			nodePortsUsed,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return desugarSenderResult{}, err
		}

		connectionsToInsert = append(connectionsToInsert, replacedConnDesugarRes.connectionsToInsert...)
		normConn = *replacedConnDesugarRes.connectionToReplace.Normal
	}

	// if sender is const (ref or literal), replace original connection with desugared and insert const and node
	if sender.Const != nil {
		if sender.Const.Value.Ref != nil {
			result, err := d.handleConstRefSender(*sender.Const.Value.Ref, nodes, scope)
			if err != nil {
				return desugarSenderResult{}, compiler.Error{
					Location: &scope.Location,
					Meta:     &sender.Meta,
				}.Wrap(err)
			}

			normConn = src.NormalConnection{
				SenderSide: []src.ConnectionSender{
					{
						PortAddr: &result,
						Meta:     sender.Meta,
					},
				},
				ReceiverSide: normConn.ReceiverSide,
			}
		}

		if sender.Const.Value.Message != nil {
			constNodePort, err := d.handleLiteralSender(*sender.Const, nodesToInsert, constsToInsert)
			if err != nil {
				return desugarSenderResult{}, err
			}

			normConn = src.NormalConnection{
				SenderSide: []src.ConnectionSender{
					{
						PortAddr: &constNodePort,
						Meta:     sender.Meta,
					},
				},
				ReceiverSide: normConn.ReceiverSide,
			}
		}
	}

	// range expression as sender
	if sender.Range != nil {
		result, err := d.handleRangeSender(
			*sender.Range,
			normConn,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return desugarSenderResult{}, err
		}
		connectionsToInsert = append(connectionsToInsert, result.connectionsToInsert...)
		normConn = result.connToReplace
	}

	return desugarSenderResult{
		connToReplace:       src.Connection{Normal: &normConn},
		connectionsToInsert: connectionsToInsert,
	}, nil
}

func (Desugarer) desugarUnnamedReceiver(
	receiver src.ConnectionPortReceiver,
	scope src.Scope,
	nodes map[string]src.Node,
) (src.ConnectionPortReceiver, *compiler.Error) {

	firstInportName, err := getFirstInportName(scope, nodes, receiver.PortAddr)
	if err != nil {
		return src.ConnectionPortReceiver{}, &compiler.Error{Err: err}
	}

	return src.ConnectionPortReceiver{
		PortAddr: src.PortAddr{
			Port: firstInportName,
			Node: receiver.PortAddr.Node,
			Idx:  receiver.PortAddr.Idx,
			Meta: receiver.PortAddr.Meta,
		},
		Meta: receiver.Meta,
	}, nil
}

func getNodeIOByPortAddr(
	scope src.Scope,
	nodes map[string]src.Node,
	portAddr *src.PortAddr,
) (src.IO, *compiler.Error) {
	node, ok := nodes[portAddr.Node]
	if !ok {
		return src.IO{}, &compiler.Error{
			Err:      fmt.Errorf("node '%s' not found", portAddr.Node),
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	entity, _, err := scope.Entity(node.EntityRef)
	if err != nil {
		return src.IO{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &portAddr.Meta,
		}
	}

	var iface src.Interface
	if entity.Kind == src.InterfaceEntity {
		iface = entity.Interface
	} else {
		iface = entity.Component.Interface
	}

	return iface.IO, nil
}

func getFirstInportName(scope src.Scope, nodes map[string]src.Node, portAddr src.PortAddr) (string, error) {
	io, err := getNodeIOByPortAddr(scope, nodes, &portAddr)
	if err != nil {
		return "", err
	}
	for inport := range io.In {
		return inport, nil
	}
	return "", errors.New("first inport not found")
}

func getFirstOutportName(scope src.Scope, nodes map[string]src.Node, portAddr src.PortAddr) (string, error) {
	io, err := getNodeIOByPortAddr(scope, nodes, &portAddr)
	if err != nil {
		return "", err
	}
	for outport := range io.Out {
		return outport, nil
	}
	return "", errors.New("first outport not found")
}

type desugarFanOutResult struct {
	receiverToReplace   src.ConnectionPortReceiver // only one (no more fan-out)
	connectionsToInsert []src.Connection
}

var fanOutCounter uint64

func (d Desugarer) desugarFanOut(
	receiverSides []src.ConnectionPortReceiver,
	nodesToInsert map[string]src.Node,
) desugarFanOutResult {
	fanOutCounter++
	nodeName := fmt.Sprintf("__fanOut__%d", fanOutCounter)

	nodesToInsert[nodeName] = src.Node{
		EntityRef: core.EntityRef{
			Name: "FanOut",
			Pkg:  "builtin",
		},
	}

	receiverToReplace := src.ConnectionPortReceiver{
		PortAddr: src.PortAddr{
			Node: nodeName,
			Port: "data",
		},
	}

	connsToInsert := make([]src.Connection, 0, len(receiverSides))
	for i, receiver := range receiverSides {
		connsToInsert = append(connsToInsert, src.Connection{
			Normal: &src.NormalConnection{
				SenderSide: []src.ConnectionSender{
					{
						PortAddr: &src.PortAddr{
							Node: nodeName,
							Port: "data",
							Idx:  compiler.Pointer(uint8(i)),
						},
					},
				},
				ReceiverSide: src.ConnectionReceiverSide{
					Receivers: []src.ConnectionPortReceiver{receiver},
				},
			},
		})
	}

	return desugarFanOutResult{
		receiverToReplace:   receiverToReplace,
		connectionsToInsert: connsToInsert,
	}
}

// Add a new atomic counter for range nodes
var rangeCounter uint64

// Add a new function to handle range senders
type handleRangeSenderResult struct {
	connToReplace       src.NormalConnection
	connectionsToInsert []src.Connection
}

func (d Desugarer) handleRangeSender(
	rangeExpr src.RangeExpr,
	normConn src.NormalConnection,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (handleRangeSenderResult, *compiler.Error) {
	rangeCounter++

	rangeNodeName := fmt.Sprintf("__range%d__", rangeCounter)
	fromConstName := fmt.Sprintf("__range%d_from__", rangeCounter)
	toConstName := fmt.Sprintf("__range%d_to__", rangeCounter)

	constsToInsert[fromConstName] = src.Const{
		TypeExpr: ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Pkg: "builtin", Name: "int"}}},
		Value:    src.ConstValue{Message: &src.MsgLiteral{Int: compiler.Pointer(int(rangeExpr.From))}},
	}
	constsToInsert[toConstName] = src.Const{
		TypeExpr: ts.Expr{Inst: &ts.InstExpr{Ref: core.EntityRef{Pkg: "builtin", Name: "int"}}},
		Value:    src.ConstValue{Message: &src.MsgLiteral{Int: compiler.Pointer(int(rangeExpr.To))}},
	}

	nodesToInsert[rangeNodeName] = src.Node{
		EntityRef: core.EntityRef{Pkg: "builtin", Name: "Range"},
	}
	nodesToInsert[fromConstName] = src.Node{
		EntityRef: core.EntityRef{Pkg: "builtin", Name: "New"},
		Directives: map[src.Directive][]string{
			"bind": {fromConstName},
		},
	}
	nodesToInsert[toConstName] = src.Node{
		EntityRef: core.EntityRef{Pkg: "builtin", Name: "New"},
		Directives: map[src.Directive][]string{
			"bind": {toConstName},
		},
	}

	connectionsToInsert := []src.Connection{
		{
			Normal: &src.NormalConnection{
				SenderSide: []src.ConnectionSender{
					{
						PortAddr: &src.PortAddr{Node: fromConstName, Port: "msg"},
					},
				},
				ReceiverSide: src.ConnectionReceiverSide{
					Receivers: []src.ConnectionPortReceiver{{
						PortAddr: src.PortAddr{Node: rangeNodeName, Port: "from"},
					}},
				},
			},
		},
		{
			Normal: &src.NormalConnection{
				SenderSide: []src.ConnectionSender{
					{
						PortAddr: &src.PortAddr{Node: toConstName, Port: "msg"},
					},
				},
				ReceiverSide: src.ConnectionReceiverSide{
					Receivers: []src.ConnectionPortReceiver{{
						PortAddr: src.PortAddr{Node: rangeNodeName, Port: "to"},
					}},
				},
			},
		},
	}

	connToReplace := src.NormalConnection{
		SenderSide: []src.ConnectionSender{
			{
				PortAddr: &src.PortAddr{Node: rangeNodeName, Port: "res"},
			},
		},
		ReceiverSide: normConn.ReceiverSide,
	}

	return handleRangeSenderResult{
		connectionsToInsert: connectionsToInsert,
		connToReplace:       connToReplace,
	}, nil
}

var fanInCounter uint64

// desugarMultipleSenders returns connections that must be used instead of given one.
// It recursevely desugars each connection before return so result is final.
func (d Desugarer) desugarMultipleSenders(
	normConn src.NormalConnection,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
	nodePortsUsed nodePortsMap,
	scope src.Scope,
	nodes map[string]src.Node,
) ([]src.Connection, error) {
	// 1. insert unique fan-in node
	fanInCounter++
	fanInNodeName := fmt.Sprintf("__fanIn__%d", fanInCounter)
	nodesToInsert[fanInNodeName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "FanIn",
		},
	}

	// 2. replace each sender with connection from sender to fan-in node
	netWithoutFanIn := make([]src.Connection, 0, len(normConn.SenderSide))
	for i, sender := range normConn.SenderSide {
		netWithoutFanIn = append(netWithoutFanIn, src.Connection{
			Normal: &src.NormalConnection{
				SenderSide: []src.ConnectionSender{sender},
				ReceiverSide: src.ConnectionReceiverSide{
					Receivers: []src.ConnectionPortReceiver{
						src.ConnectionPortReceiver{
							PortAddr: src.PortAddr{
								Node: fanInNodeName,
								Port: "data",
								Idx:  compiler.Pointer(uint8(i)),
							},
						},
					},
				},
			},
		})
	}

	// 3. insert new connection from fan-in node to original receiver
	netWithoutFanIn = append(netWithoutFanIn, src.Connection{
		Normal: &src.NormalConnection{
			SenderSide: []src.ConnectionSender{
				{
					PortAddr: &src.PortAddr{
						Node: fanInNodeName,
						Port: "res",
					},
				},
			},
			ReceiverSide: src.ConnectionReceiverSide{
				Receivers: normConn.ReceiverSide.Receivers,
			},
		},
	})

	// 4. desugar each connection (original senders and receivers might need it)
	desugaredConnections, err := d.desugarConnections(
		netWithoutFanIn,
		nodePortsUsed,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
	)
	if err != nil {
		return nil, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &normConn.Meta,
		}
	}

	return desugaredConnections, nil
}
