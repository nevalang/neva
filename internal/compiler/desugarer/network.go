package desugarer

import (
	"errors"
	"fmt"
	"maps"
	"sync/atomic"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

type handleNetResult struct {
	desugaredConnections []src.Connection        // desugared network
	constsToInsert       map[string]src.ConstDef // constants that needs to be inserted in to make desugared network work
	nodesToInsert        map[string]src.Node     // nodes that needs to be inserted in to make desugared network work
	nodesPortsUsed       nodePortsMap            // to find unused to create virtual del connections
}

func (d Desugarer) handleNetwork(
	net []src.Connection,
	nodes map[string]src.Node,
	scope src.Scope,
) (handleNetResult, *compiler.Error) {
	// code smell: mix of mutable and immutable styles (connections/nodes-consts)
	desugaredConnections := make([]src.Connection, 0, len(net))
	nodesToInsert := map[string]src.Node{}
	constsToInsert := map[string]src.ConstDef{}
	nodesPortsUsed := newNodePortsMap()

	for _, conn := range net {
		result, err := d.desugarConnection(
			conn,
			nodesPortsUsed,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return handleNetResult{}, err
		}

		desugaredConnections = append(desugaredConnections, result.connectionToReplace)
		desugaredConnections = append(desugaredConnections, result.connectionsToInsert...)
	}

	result, err := d.networkFinalProcessing(desugaredConnections, nodesToInsert)
	if err != nil {
		return handleNetResult{}, &compiler.Error{} // todo
	}
	desugaredConnections = result.FinalNetwork

	return handleNetResult{
		desugaredConnections: desugaredConnections,
		nodesPortsUsed:       nodesPortsUsed,
		constsToInsert:       constsToInsert,
		nodesToInsert:        nodesToInsert,
	}, nil
}

type NetFinalProcessingResult struct {
	FinalNetwork []src.Connection
}

var fanInCounter *atomic.Uint64 = &atomic.Uint64{}

// networkFinalProcessing does the following:
//   - inserts fan-in node in-between everywhere needed: generates necessary connections and inserts corresponding nodes
//   - ...
//
// It's important that given conns must be already desugared separately.
func (d Desugarer) networkFinalProcessing(
	conns []src.Connection,
	nodesToInsert map[string]src.Node,
) (NetFinalProcessingResult, error) {
	finalNet := []src.Connection{}
	receiverSendersMap := make(map[src.PortAddr][]src.PortAddr) // receiver -> senders

	for _, conn := range conns {
		if conn.ArrayBypass != nil {
			finalNet = append(finalNet, conn) // nothing to desugar
			continue
		}

		senderAddr := *conn.Normal.SenderSide.PortAddr
		receiverAddr := conn.Normal.ReceiverSide.Receivers[0].PortAddr
		receiverSendersMap[receiverAddr] = append(receiverSendersMap[receiverAddr], senderAddr)
	}

	for receiver, senders := range receiverSendersMap {
		if len(senders) < 2 { // keep non fan-in connections as-is
			finalNet = append(finalNet, src.Connection{
				Normal: &src.NormalConnection{
					ReceiverSide: src.ConnectionReceiverSide{
						Receivers: []src.ConnectionReceiver{{PortAddr: receiver}},
					},
					SenderSide: src.ConnectionSenderSide{PortAddr: &senders[0]},
				},
			})
			continue
		}

		// ---
		// for each connection with >1 senders we need to
		// 1. create unique fan-in node
		// 2. replace existing receiver with created fan-in
		// 3. create new connection with fan-in as sender and original receiver as receiver
		// ---

		// 1. create unique fan-in node
		fanInNodeName := fmt.Sprintf("__fanIn__%d", fanInCounter.Add(1))
		nodesToInsert[fanInNodeName] = src.Node{
			EntityRef: core.EntityRef{
				Pkg:  "builtin",
				Name: "FanIn",
			},
		}

		for i, sender := range senders {
			s := sender

			// 2. replace existing receiver with created fan-in
			finalNet = append(finalNet, src.Connection{
				Normal: &src.NormalConnection{
					SenderSide: src.ConnectionSenderSide{
						PortAddr: &s,
					},
					ReceiverSide: src.ConnectionReceiverSide{
						Receivers: []src.ConnectionReceiver{
							{
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

			// 3. create new connection with fan-in as sender and original receiver as receiver
			finalNet = append(finalNet, src.Connection{
				Normal: &src.NormalConnection{
					SenderSide: src.ConnectionSenderSide{
						PortAddr: &src.PortAddr{
							Node: fanInNodeName,
							Port: "res",
						},
					},
					ReceiverSide: src.ConnectionReceiverSide{
						Receivers: []src.ConnectionReceiver{
							{PortAddr: receiver},
						},
					},
				},
			})
		}
	}

	return NetFinalProcessingResult{
		FinalNetwork: finalNet,
	}, nil
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
	constsToInsert map[string]src.ConstDef,
) (desugarConnectionResult, *compiler.Error) {
	// "array bypass" connection - nothing to desugar, just mark as used and return as-is
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

	// further we only handle normal connections

	// mark as used and handle unnamed port if needed
	if conn.Normal.SenderSide.PortAddr != nil {
		if conn.Normal.SenderSide.PortAddr.Port == "" {
			firstOutportName, err := getFirstOutportName(scope, nodes, *conn.Normal.SenderSide.PortAddr)
			if err != nil {
				return desugarConnectionResult{}, &compiler.Error{Err: err}
			}

			conn = src.Connection{
				Normal: &src.NormalConnection{
					SenderSide: src.ConnectionSenderSide{
						PortAddr: &src.PortAddr{
							Port: firstOutportName,
							Node: conn.Normal.SenderSide.PortAddr.Node,
							Idx:  conn.Normal.SenderSide.PortAddr.Idx,
							Meta: conn.Normal.SenderSide.PortAddr.Meta,
						},
						Selectors: conn.Normal.SenderSide.Selectors,
						Meta:      conn.Normal.SenderSide.Meta,
					},
					ReceiverSide: conn.Normal.ReceiverSide,
				},
				Meta: conn.Meta,
			}
		}

		// mark as used
		nodePortsUsed.set(
			conn.Normal.SenderSide.PortAddr.Node,
			conn.Normal.SenderSide.PortAddr.Port,
		)
	}

	connectionsToInsert := []src.Connection{}

	// if conn has selectors, desugar them, replace original connection and insert what's needed
	if len(conn.Normal.SenderSide.Selectors) != 0 {
		result, err := d.desugarStructSelectors(*conn.Normal)
		if err != nil {
			return desugarConnectionResult{}, compiler.Error{
				Err:      errors.New("Cannot desugar struct selectors"),
				Location: &scope.Location,
				Meta:     &conn.Meta,
			}.Wrap(err)
		}

		nodesToInsert[result.nodeToInsertName] = result.nodeToInsert
		constsToInsert[result.constToInsertName] = result.constToInsert

		// generated connection might need desugaring itself
		connToInsertDesugarRes, err := d.desugarConnection(
			result.connToInsert,
			nodePortsUsed,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return desugarConnectionResult{}, err
		}

		connectionsToInsert = append(connectionsToInsert, connToInsertDesugarRes.connectionToReplace)
		connectionsToInsert = append(connectionsToInsert, connToInsertDesugarRes.connectionsToInsert...)

		// connection that replaces original one might need desugaring itself
		replacedConnDesugarRes, err := d.desugarConnection(
			result.connToReplace,
			nodePortsUsed,
			scope,
			nodes,
			nodesToInsert,
			constsToInsert,
		)
		if err != nil {
			return desugarConnectionResult{}, err
		}

		connectionsToInsert = append(connectionsToInsert, replacedConnDesugarRes.connectionsToInsert...)

		conn = replacedConnDesugarRes.connectionToReplace
	}

	// if sender is const (ref or literal), replace original connection with desugared and insert const and node
	if conn.Normal.SenderSide.Const != nil {
		if conn.Normal.SenderSide.Const.Value.Ref != nil {
			result, err := d.handleConstRefSender(conn, scope)
			if err != nil {
				return desugarConnectionResult{}, err
			}
			nodesToInsert[result.nodeToInsertName] = result.nodeToInsert
			conn = result.connToReplace
		} else if conn.Normal.SenderSide.Const.Value.Message != nil {
			result, err := d.handleLiteralSender(conn)
			if err != nil {
				return desugarConnectionResult{}, err
			}
			constsToInsert[result.constName] = *conn.Normal.SenderSide.Const
			nodesToInsert[result.nodeToInsertName] = result.nodeToInsert
			conn = result.connToReplace
		}
	}

	desugaredReceivers := make([]src.ConnectionReceiver, 0, len(conn.Normal.ReceiverSide.Receivers))

	// desugar unnamed receivers if needed and replace them with named ones
	for _, receiver := range conn.Normal.ReceiverSide.Receivers {
		if receiver.PortAddr.Port != "" {
			desugaredReceivers = append(desugaredReceivers, receiver)
			continue
		}

		firstInportName, err := getFirstInportName(scope, nodes, receiver.PortAddr)
		if err != nil {
			return desugarConnectionResult{}, &compiler.Error{Err: err}
		}

		desugaredReceivers = append(desugaredReceivers, src.ConnectionReceiver{
			PortAddr: src.PortAddr{
				Port: firstInportName,
				Node: receiver.PortAddr.Node,
				Idx:  receiver.PortAddr.Idx,
				Meta: receiver.PortAddr.Meta,
			},
			Meta: receiver.Meta,
		})
	}

	// it's possible to have connection with both normal receivers and deferred connections so we handle both

	if conn.Normal.ReceiverSide.DeferredConnections != nil {
		result, err := d.desugarDeferredConnections(
			*conn.Normal,
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
		result := d.desugarFanOut(desugaredReceivers)
		nodesToInsert[result.nodeToInsertName] = result.nodeToInsert
		desugaredReceivers = []src.ConnectionReceiver{result.receiverToReplace} // replace all existing receivers with single one
		connectionsToInsert = append(connectionsToInsert, result.connectionsToInsert...)
	}

	return desugarConnectionResult{
		connectionToReplace: src.Connection{
			Normal: &src.NormalConnection{
				SenderSide: conn.Normal.SenderSide,
				ReceiverSide: src.ConnectionReceiverSide{
					Receivers: desugaredReceivers,
				},
			},
			Meta: conn.Meta,
		},
		connectionsToInsert: connectionsToInsert,
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
	nodeToInsertName    string
	nodeToInsert        src.Node
	receiverToReplace   src.ConnectionReceiver // only one (no more fan-out)
	connectionsToInsert []src.Connection
}

var fanOutCounter atomic.Uint64

func (d Desugarer) desugarFanOut(receiverSides []src.ConnectionReceiver) desugarFanOutResult {
	counter := fanOutCounter.Load()
	fanOutCounter.Store(counter + 1)
	nodeName := fmt.Sprintf("__fanOut__%d", counter)

	node := src.Node{
		EntityRef: core.EntityRef{
			Name: "FanOut",
			Pkg:  "builtin",
		},
	}

	receiverToReplace := src.ConnectionReceiver{
		PortAddr: src.PortAddr{
			Node: nodeName,
			Port: "data",
		},
	}

	connsToInsert := make([]src.Connection, 0, len(receiverSides))
	for i, receiver := range receiverSides {
		connsToInsert = append(connsToInsert, src.Connection{
			Normal: &src.NormalConnection{
				SenderSide: src.ConnectionSenderSide{
					PortAddr: &src.PortAddr{
						Node: nodeName,
						Port: "data",
						Idx:  compiler.Pointer(uint8(i)),
					},
				},
				ReceiverSide: src.ConnectionReceiverSide{
					Receivers: []src.ConnectionReceiver{receiver},
				},
			},
		})
	}

	return desugarFanOutResult{
		nodeToInsertName:    nodeName,
		nodeToInsert:        node,
		receiverToReplace:   receiverToReplace,
		connectionsToInsert: connsToInsert,
	}
}
