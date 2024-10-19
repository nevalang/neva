package desugarer

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	"github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

type handleNetworkResult struct {
	desugaredConnections []src.Connection
	constsToInsert       map[string]src.Const
	nodesToInsert        map[string]src.Node
	nodesPortsUsed       nodeOutportsUsed
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
	nodePortsUsed nodeOutportsUsed,
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
			nil,
		)
		if err != nil {
			return nil, err
		}

		if result.replace != nil {
			desugaredConnections = append(desugaredConnections, *result.replace)
		}

		desugaredConnections = append(desugaredConnections, result.insert...)
	}

	return desugaredConnections, nil
}

type desugarConnectionResult struct {
	replace *src.Connection
	insert  []src.Connection
}

func (d Desugarer) desugarConnection(
	conn src.Connection,
	nodePortsUsed nodeOutportsUsed,
	scope src.Scope,
	nodes map[string]src.Node,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
	prevChainLink *src.ConnectionSender,
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
		return desugarConnectionResult{replace: &conn}, nil
	}

	return d.desugarNormalConnection(
		*conn.Normal,
		nodePortsUsed,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
		prevChainLink,
	)
}

func (d Desugarer) desugarNormalConnection(
	normConn src.NormalConnection,
	nodePortsUsed nodeOutportsUsed,
	scope src.Scope,
	nodes map[string]src.Node,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
	prevChainLink *src.ConnectionSender,
) (desugarConnectionResult, *compiler.Error) {
	if len(normConn.SenderSide) > 1 {
		result, err := d.desugarFanIn(
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
			insert: result,
		}, nil
	}

	desugarSenderResult, err := d.desugarSingleSender(
		normConn,
		scope,
		nodes,
		nodePortsUsed,
		nodesToInsert,
		constsToInsert,
		prevChainLink,
	)
	if err != nil {
		return desugarConnectionResult{}, compiler.Error{
			Location: &scope.Location,
			Meta:     &normConn.SenderSide[0].Meta,
		}.Wrap(err)
	}

	normConn = *desugarSenderResult.replace.Normal
	insert := desugarSenderResult.insert

	if len(normConn.ReceiverSide) > 1 {
		result, err := d.desugarFanOut(
			normConn,
			nodesToInsert,
			constsToInsert,
			nodePortsUsed,
			scope,
			nodes,
		)
		if err != nil {
			return desugarConnectionResult{}, err
		}

		return desugarConnectionResult{
			replace: &result.replace,
			insert:  append(insert, result.insert...),
		}, nil
	}

	desugarReceiverResult, err := d.desugarSingleReceiver(
		normConn,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
		nodePortsUsed,
	)
	if err != nil {
		return desugarConnectionResult{}, err
	}

	return desugarConnectionResult{
		replace: &desugarReceiverResult.replace,
		insert:  append(insert, desugarReceiverResult.insert...),
	}, nil
}

type desugarReceiverResult struct {
	replace src.Connection
	insert  []src.Connection
}

func (d Desugarer) desugarSingleReceiver(
	normConn src.NormalConnection,
	scope src.Scope,
	nodes map[string]src.Node,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
	nodePortsUsed nodeOutportsUsed,
) (desugarReceiverResult, *compiler.Error) {
	receiver := normConn.ReceiverSide[0]

	if receiver.PortAddr != nil {
		if receiver.PortAddr.Port != "" {
			return desugarReceiverResult{
				replace: src.Connection{
					Normal: &src.NormalConnection{
						SenderSide:   normConn.SenderSide,
						ReceiverSide: []src.ConnectionReceiver{receiver},
					},
				},
				insert: []src.Connection{},
			}, nil
		}

		firstInportName, err := getFirstInportName(scope, nodes, *receiver.PortAddr)
		if err != nil {
			return desugarReceiverResult{}, &compiler.Error{
				Err:      err,
				Location: &scope.Location,
				Meta:     &receiver.Meta,
			}
		}

		return desugarReceiverResult{
			replace: src.Connection{
				Normal: &src.NormalConnection{
					SenderSide: normConn.SenderSide,
					ReceiverSide: []src.ConnectionReceiver{
						{
							PortAddr: &src.PortAddr{
								Port: firstInportName,
								Node: receiver.PortAddr.Node,
								Idx:  receiver.PortAddr.Idx,
								Meta: receiver.PortAddr.Meta,
							},
						},
					},
				},
			},
		}, nil
	}

	if receiver.DeferredConnection != nil {
		result, err := d.desugarDeferredConnection(
			normConn,
			scope,
			constsToInsert,
			nodesToInsert,
			nodePortsUsed,
			nodes,
		)
		if err != nil {
			return desugarReceiverResult{}, err
		}

		return desugarReceiverResult(result), nil
	}

	desugarChainResult, err := d.desugarChainedConnection(
		receiver,
		scope,
		nodes,
		nodePortsUsed,
		nodesToInsert,
		constsToInsert,
		normConn,
	)
	if err != nil {
		return desugarReceiverResult{}, err
	}

	return desugarReceiverResult{
		replace: *desugarChainResult.replace,
		insert:  desugarChainResult.insert,
	}, nil
}

func (d Desugarer) desugarChainedConnection(
	receiver src.ConnectionReceiver,
	scope src.Scope,
	nodes map[string]src.Node,
	nodePortsUsed nodeOutportsUsed,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
	normConn src.NormalConnection,
) (desugarConnectionResult, *compiler.Error) {
	chainedConn := *receiver.ChainedConnection
	chainHead := chainedConn.Normal.SenderSide[0] // chain head is always single sender

	// it's only possible to find receiver port before desugaring
	var chainHeadPort string
	if chainHead.Range != nil {
		chainHeadPort = "sig"
	} else if chainHead.PortAddr != nil {
		var firstInportName = chainHead.PortAddr.Port
		if chainHead.PortAddr.Port == "" {
			var err error
			firstInportName, err = getFirstInportName(scope, nodes, *chainHead.PortAddr)
			if err != nil {
				return desugarConnectionResult{}, &compiler.Error{Err: err}
			}
		}
		chainHeadPort = firstInportName
	}

	desugarChainResult, err := d.desugarConnection(
		chainedConn,
		nodePortsUsed,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
		&normConn.SenderSide[0],
	)
	if err != nil {
		return desugarConnectionResult{}, err
	}

	desugaredHead := desugarChainResult.replace.Normal.SenderSide[0]

	replace := src.Connection{
		Normal: &src.NormalConnection{
			SenderSide: normConn.SenderSide,
			ReceiverSide: []src.ConnectionReceiver{
				{
					PortAddr: &src.PortAddr{
						Node: desugaredHead.PortAddr.Node,
						Port: chainHeadPort,
						Idx:  desugaredHead.PortAddr.Idx,
						Meta: chainHead.Meta,
					},
				},
			},
		},
	}

	// we need to insert both: replace and insert, example:
	// input = a -> b -> c -> d
	// sender = a
	// chain = b -> c -> d
	// ---
	// replace = b -> c
	// insert = c -> d
	// (and replace existing one with `a -> b`)
	insert := append([]src.Connection{}, desugarChainResult.insert...)
	insert = append(insert, *desugarChainResult.replace)

	return desugarConnectionResult{
		replace: &replace,
		insert:  insert,
	}, nil
}

type desugarDeferredConnectionsResult struct {
	replace src.Connection
	insert  []src.Connection
}

var virtualLocksCounter uint64

func (d Desugarer) desugarDeferredConnection(
	normConn src.NormalConnection,
	scope src.Scope,
	constsToInsert map[string]src.Const,
	nodesToInsert map[string]src.Node,
	nodesPortsUsed nodeOutportsUsed,
	nodes map[string]src.Node,
) (desugarDeferredConnectionsResult, *compiler.Error) {
	deferredConnection := *normConn.ReceiverSide[0].DeferredConnection

	desugarDeferredConnResult, err := d.desugarConnection(
		deferredConnection,
		nodesPortsUsed,
		scope,
		nodes,
		nodesToInsert,
		constsToInsert,
		nil,
	)
	if err != nil {
		return desugarDeferredConnectionsResult{}, err
	}

	deferredConnection = *desugarDeferredConnResult.replace
	connsToInsert := desugarDeferredConnResult.insert

	// 1) create lock node
	virtualLocksCounter++
	lockNodeName := fmt.Sprintf("__lock__%d", virtualLocksCounter)
	nodesToInsert[lockNodeName] = src.Node{
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "Lock",
		},
		TypeArgs: []typesystem.Expr{
			ts.Expr{
				Inst: &typesystem.InstExpr{
					Ref: core.EntityRef{Pkg: "builtin", Name: "any"},
				},
			},
		},
	}

	// 2) connect original sender to lock receiver
	replace := src.Connection{
		Normal: &src.NormalConnection{
			SenderSide: normConn.SenderSide,
			ReceiverSide: []src.ConnectionReceiver{
				{
					PortAddr: &src.PortAddr{
						Node: lockNodeName,
						Port: "sig",
					},
				},
			},
		},
	}

	connsToInsert = append(
		// 3) connect deferred sender to lock data
		connsToInsert,
		src.Connection{
			Normal: &src.NormalConnection{
				SenderSide: deferredConnection.Normal.SenderSide,
				ReceiverSide: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{
							Node: lockNodeName,
							Port: "data",
						},
					},
				},
			},
		},
		// 4) create connection from lock:data to receiver-side of deferred connection
		src.Connection{
			Normal: &src.NormalConnection{
				SenderSide: []src.ConnectionSender{
					{
						PortAddr: &src.PortAddr{
							Node: lockNodeName,
							Port: "data",
						},
					},
				},
				ReceiverSide: deferredConnection.Normal.ReceiverSide,
			},
		},
	)

	return desugarDeferredConnectionsResult{
		replace: replace,
		insert:  connsToInsert,
	}, nil
}

type desugarSenderResult struct {
	replace src.Connection
	insert  []src.Connection
}

func (d Desugarer) desugarSingleSender(
	normConn src.NormalConnection,
	scope src.Scope,
	nodes map[string]src.Node,
	usedNodeOutports nodeOutportsUsed,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
	prevChainLink *src.ConnectionSender,
) (desugarSenderResult, *compiler.Error) {
	sender := normConn.SenderSide[0]

	// mark outport as used and desugar unnamed port if needed
	if sender.PortAddr != nil {
		portName := sender.PortAddr.Port
		if sender.PortAddr.Port == "" {
			firstOutportName, err := getFirstOutportName(scope, nodes, *sender.PortAddr)
			if err != nil {
				return desugarSenderResult{}, &compiler.Error{Err: err}
			}
			portName = firstOutportName
			normConn.SenderSide = []src.ConnectionSender{
				{
					PortAddr: &src.PortAddr{
						Port: portName,
						Node: sender.PortAddr.Node,
						Idx:  sender.PortAddr.Idx,
						Meta: sender.PortAddr.Meta,
					},
				},
			}
		}
		usedNodeOutports.set(
			sender.PortAddr.Node,
			portName,
		)
		return desugarSenderResult{
			replace: src.Connection{Normal: &normConn},
			insert:  nil,
		}, nil
	}

	// if conn has selectors, desugar them, replace original connection and insert what's needed
	if len(sender.Selectors) != 0 {
		result, err := d.desugarStructSelectors(
			*prevChainLink,
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
			result.insert,
			usedNodeOutports,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
			nil,
		)
		if err != nil {
			return desugarSenderResult{}, err
		}

		// connection that replaces original one might need desugaring itself
		replacedConnDesugarRes, err := d.desugarConnection(
			result.connToReplace,
			usedNodeOutports,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
			nil,
		)
		if err != nil {
			return desugarSenderResult{}, err
		}

		insert := []src.Connection{}
		insert = append(insert, *connToInsertDesugarRes.replace)
		insert = append(insert, connToInsertDesugarRes.insert...)

		return desugarSenderResult{
			replace: src.Connection{Normal: replacedConnDesugarRes.replace.Normal},
			insert:  append(insert, replacedConnDesugarRes.insert...),
		}, nil
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

		return desugarSenderResult{
			replace: src.Connection{Normal: &normConn},
			insert:  nil,
		}, nil
	}

	result, err := d.desugarRangeSender(
		*sender.Range,
		normConn,
		nodesToInsert,
		constsToInsert,
	)
	if err != nil {
		return desugarSenderResult{}, err
	}

	return desugarSenderResult{
		replace: src.Connection{Normal: &result.replace},
		insert:  result.insert,
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
	replace src.Connection   // original sender -> fanOut receiver
	insert  []src.Connection // fanOut sender -> original receivers
}

var fanOutCounter uint64

func (d Desugarer) desugarFanOut(
	normConn src.NormalConnection,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
	nodePortsUsed nodeOutportsUsed,
	scope src.Scope,
	nodes map[string]src.Node,
) (desugarFanOutResult, *compiler.Error) {
	fanOutCounter++
	nodeName := fmt.Sprintf("__fanOut__%d", fanOutCounter)

	nodesToInsert[nodeName] = src.Node{
		EntityRef: core.EntityRef{
			Name: "FanOut",
			Pkg:  "builtin",
		},
	}

	receiverToReplace := src.ConnectionReceiver{
		PortAddr: &src.PortAddr{
			Node: nodeName,
			Port: "data",
		},
	}

	insert := make([]src.Connection, 0, len(normConn.ReceiverSide))
	for i, receiver := range normConn.ReceiverSide {
		conn := src.Connection{
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
				ReceiverSide: []src.ConnectionReceiver{receiver},
			},
		}

		desugarConnRes, err := d.desugarConnection(
			conn,
			nodePortsUsed,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
			nil,
		)
		if err != nil {
			return desugarFanOutResult{}, err
		}

		insert = append(insert, *desugarConnRes.replace)
		insert = append(insert, desugarConnRes.insert...)
	}

	return desugarFanOutResult{
		replace: src.Connection{
			Normal: &src.NormalConnection{
				SenderSide:   normConn.SenderSide, // senders must be desugared
				ReceiverSide: []src.ConnectionReceiver{receiverToReplace},
			},
		},
		insert: insert,
	}, nil
}

// Add a new atomic counter for range nodes
var rangeCounter uint64

// Add a new function to handle range senders
type handleRangeSenderResult struct {
	replace src.NormalConnection
	insert  []src.Connection
}

// desugarRangeSender desugars `from..to -> XXX` part.
// It does not create connection to range:sig,
// it's done in chained connection desugaring.
func (d Desugarer) desugarRangeSender(
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

	replace := src.NormalConnection{
		SenderSide: []src.ConnectionSender{
			{
				PortAddr: &src.PortAddr{Node: rangeNodeName, Port: "res"},
			},
		},
		ReceiverSide: normConn.ReceiverSide,
	}

	insert := []src.Connection{
		// $from -> range:from
		{
			Normal: &src.NormalConnection{
				SenderSide: []src.ConnectionSender{
					{
						PortAddr: &src.PortAddr{Node: fromConstName, Port: "msg"},
					},
				},
				ReceiverSide: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{Node: rangeNodeName, Port: "from"},
					},
				},
			},
		},
		// $to -> range:to
		{
			Normal: &src.NormalConnection{
				SenderSide: []src.ConnectionSender{
					{
						PortAddr: &src.PortAddr{Node: toConstName, Port: "msg"},
					},
				},
				ReceiverSide: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{Node: rangeNodeName, Port: "to"},
					},
				},
			},
		},
	}

	return handleRangeSenderResult{
		insert:  insert,
		replace: replace,
	}, nil
}

var fanInCounter uint64

// desugarFanIn returns connections that must be used instead of given one.
// It recursevely desugars each connection before return so result is final.
func (d Desugarer) desugarFanIn(
	normConn src.NormalConnection,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
	nodePortsUsed nodeOutportsUsed,
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

	// 2. connection each sender with fan-in node
	netWithoutFanIn := make([]src.Connection, 0, len(normConn.SenderSide))
	for i, sender := range normConn.SenderSide {
		netWithoutFanIn = append(netWithoutFanIn, src.Connection{
			Normal: &src.NormalConnection{
				SenderSide: []src.ConnectionSender{sender},
				ReceiverSide: []src.ConnectionReceiver{
					{
						PortAddr: &src.PortAddr{
							Node: fanInNodeName,
							Port: "data",
							Idx:  compiler.Pointer(uint8(i)),
						},
					},
				},
			},
		})
	}

	// 3. insert new connection from fan-in to original receivers
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
			ReceiverSide: normConn.ReceiverSide,
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
