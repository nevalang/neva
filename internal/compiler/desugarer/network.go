package desugarer

import (
	"errors"
	"maps"
	"slices"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

type handleNetResult struct {
	desugaredConnections []src.Connection     // desugared network
	virtualConstants     map[string]src.Const // constants that needs to be inserted in to make desugared network work
	virtualNodes         map[string]src.Node  // nodes that needs to be inserted in to make desugared network work
	usedNodePorts        nodePortsMap         // to find unused to create virtual del connections
}

func (d Desugarer) handleNetwork(
	net []src.Connection,
	nodes map[string]src.Node,
	scope src.Scope,
) (handleNetResult, *compiler.Error) {
	desugaredConns := make([]src.Connection, 0, len(net))
	nodesToInsert := map[string]src.Node{}
	constsToInsert := map[string]src.Const{}
	usedNodePorts := newNodePortsMap()

	for _, conn := range net {
		if conn.ArrayBypass != nil {
			usedNodePorts.set(
				conn.ArrayBypass.SenderOutport.Node,
				conn.ArrayBypass.SenderOutport.Port,
			)
			usedNodePorts.set(
				conn.ArrayBypass.ReceiverInport.Node,
				conn.ArrayBypass.ReceiverInport.Port,
			)
			desugaredConns = append(desugaredConns, conn)
			continue
		}

		if conn.Normal.SenderSide.PortAddr != nil {
			if conn.Normal.SenderSide.PortAddr.Port != "" {
				usedNodePorts.set(
					conn.Normal.SenderSide.PortAddr.Node,
					conn.Normal.SenderSide.PortAddr.Port,
				)
				continue
			}

			found, err := getFirstOutPortName(scope, nodes, *conn.Normal.SenderSide.PortAddr)
			if err != nil {
				return handleNetResult{}, &compiler.Error{
					Err: err,
				}
			}

			conn = src.Connection{
				Normal: &src.NormalConnection{
					SenderSide: src.ConnectionSenderSide{
						PortAddr: &src.PortAddr{
							Port: found,
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

			usedNodePorts.set(
				conn.Normal.SenderSide.PortAddr.Node,
				conn.Normal.SenderSide.PortAddr.Port,
			)
		}

		if len(conn.Normal.SenderSide.Selectors) != 0 {
			result, err := d.desugarStructSelectors(*conn.Normal)
			if err != nil {
				return handleNetResult{}, compiler.Error{
					Err:      errors.New("Cannot desugar struct selectors"),
					Location: &scope.Location,
					Meta:     &conn.Meta,
				}.Wrap(err)
			}
			nodesToInsert[result.nodeToInsertName] = result.nodeToInsert
			constsToInsert[result.constToInsertName] = result.constToInsert
			conn = result.connToReplace

			// FIXME result.connToInsert needs desugaring itself! Because of portless conn
			desugaredConns = append(desugaredConns, result.connToInsert)
		}

		if conn.Normal.SenderSide.Const != nil {
			if conn.Normal.SenderSide.Const.Ref != nil {
				result, err := d.handleConstRefSender(conn, scope)
				if err != nil {
					return handleNetResult{}, err
				}
				nodesToInsert[result.emitterNodeName] = result.emitterNode
				conn = result.connectionWithoutConstSender
			} else if conn.Normal.SenderSide.Const.Message != nil {
				result, err := d.handleLiteralSender(conn)
				if err != nil {
					return handleNetResult{}, err
				}
				constsToInsert[result.constName] = *conn.Normal.SenderSide.Const
				nodesToInsert[result.emitterNodeName] = result.emitterNode
				conn = result.connectionWithoutConstSender
			}
		}

		if len(conn.Normal.ReceiverSide.DeferredConnections) == 0 {
			desugaredReceivers := slices.Clone(conn.Normal.ReceiverSide.Receivers)

			for i, receiver := range conn.Normal.ReceiverSide.Receivers {
				if receiver.PortAddr.Port != "" {
					usedNodePorts.set(receiver.PortAddr.Node, receiver.PortAddr.Port)
					continue
				}

				found, err := getFirstInPortName(scope, nodes, receiver.PortAddr)
				if err != nil {
					return handleNetResult{}, &compiler.Error{Err: err}
				}

				desugaredReceivers[i] = src.ConnectionReceiver{
					PortAddr: src.PortAddr{
						Port: found,
						Node: receiver.PortAddr.Node,
						Idx:  receiver.PortAddr.Idx,
						Meta: receiver.PortAddr.Meta,
					},
					Meta: receiver.Meta,
				}

				usedNodePorts.set(receiver.PortAddr.Node, receiver.PortAddr.Port)
			}

			conn = src.Connection{
				Normal: &src.NormalConnection{
					SenderSide: conn.Normal.SenderSide,
					ReceiverSide: src.ConnectionReceiverSide{
						Receivers: desugaredReceivers,
					},
				},
				Meta: conn.Meta,
			}

			desugaredConns = append(desugaredConns, conn)
			continue
		}

		deferredConnsResult, err := d.handleDeferredConnections(
			*conn.Normal,
			nodes,
			scope,
		)
		if err != nil {
			return handleNetResult{}, err
		}

		maps.Copy(usedNodePorts.m, deferredConnsResult.usedNodesPorts.m)
		maps.Copy(constsToInsert, deferredConnsResult.virtualConstants)
		maps.Copy(nodesToInsert, deferredConnsResult.virtualNodes)

		desugaredConns = append(
			desugaredConns,
			deferredConnsResult.desugaredConnections...,
		)
	}

	return handleNetResult{
		desugaredConnections: desugaredConns,
		usedNodePorts:        usedNodePorts,
		virtualConstants:     constsToInsert,
		virtualNodes:         nodesToInsert,
	}, nil
}

func getNodeIOByPortAddr(
	scope src.Scope,
	nodes map[string]src.Node,
	portAddr *src.PortAddr,
) (src.IO, *compiler.Error) {
	entity, _, err := scope.Entity(nodes[portAddr.Node].EntityRef)
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

func getFirstInPortName(scope src.Scope, nodes map[string]src.Node, portAddr src.PortAddr) (string, error) {
	io, err := getNodeIOByPortAddr(scope, nodes, &portAddr)
	if err != nil {
		return "", err
	}
	for inport := range io.In {
		return inport, nil
	}
	return "", errors.New("first inport not found")
}

func getFirstOutPortName(scope src.Scope, nodes map[string]src.Node, portAddr src.PortAddr) (string, error) {
	io, err := getNodeIOByPortAddr(scope, nodes, &portAddr)
	if err != nil {
		return "", err
	}
	for outport := range io.Out {
		return outport, nil
	}
	return "", errors.New("first outport not found")
}
