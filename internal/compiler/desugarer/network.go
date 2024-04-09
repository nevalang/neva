package desugarer

import (
	"errors"
	"maps"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
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
	nodesToInsert := map[string]src.Node{}
	desugaredConns := make([]src.Connection, 0, len(net))
	usedNodePorts := newNodePortsMap()
	constantsToInsert := map[string]src.Const{}

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
			if conn.Normal.SenderSide.PortAddr.Port == "" { // if port is not specified, use first available
				io, err := getNodeIOByPortAddr(scope, nodes, conn.Normal.SenderSide.PortAddr)
				if err != nil {
					return handleNetResult{}, err
				}
				for outport := range io.Out {
					conn = src.Connection{
						Normal: &src.NormalConnection{
							SenderSide: src.ConnectionSenderSide{
								PortAddr: &src.PortAddr{
									Node: conn.Normal.SenderSide.PortAddr.Node,
									Port: outport, // <- substituted
									Idx:  conn.Normal.SenderSide.PortAddr.Idx,
									Meta: core.Meta{},
								},
								Selectors: conn.Normal.SenderSide.Selectors,
								Meta:      conn.Normal.SenderSide.Meta,
							},
							ReceiverSide: conn.Normal.ReceiverSide,
						},
						Meta: conn.Meta,
					}
					conn.Normal.SenderSide.PortAddr.Port = outport
					break
				}
			}

			usedNodePorts.set(
				conn.Normal.SenderSide.PortAddr.Node,
				conn.Normal.SenderSide.PortAddr.Port,
			)
		}

		// TODO possible option is to check whether we do not need to desugar anything here
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
			constantsToInsert[result.constToInsertName] = result.constToInsert
			conn = result.connToReplace
			desugaredConns = append(desugaredConns, result.connToInsert)
		}

		if conn.Normal.SenderSide.Const != nil { //nolint:nestif
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
				constantsToInsert[result.constName] = *conn.Normal.SenderSide.Const
				nodesToInsert[result.emitterNodeName] = result.emitterNode
				conn = result.connectionWithoutConstSender
			}
		}

		if len(conn.Normal.ReceiverSide.DeferredConnections) == 0 {
			for i, receiver := range conn.Normal.ReceiverSide.Receivers {
				if receiver.PortAddr.Port == "" {
					io, err := getNodeIOByPortAddr(scope, nodes, &receiver.PortAddr)
					if err != nil {
						return handleNetResult{}, err
					}
					for inport := range io.In {
						receiver.PortAddr.Port = inport
						break
					}
					conn.Normal.ReceiverSide.Receivers[i].PortAddr.Port = receiver.PortAddr.Port // <- substituted (mutation)
				}

				usedNodePorts.set(receiver.PortAddr.Node, receiver.PortAddr.Port)
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
		maps.Copy(constantsToInsert, deferredConnsResult.virtualConstants)
		maps.Copy(nodesToInsert, deferredConnsResult.virtualNodes)

		desugaredConns = append(
			desugaredConns,
			deferredConnsResult.desugaredConnections...,
		)
	}

	return handleNetResult{
		desugaredConnections: desugaredConns,
		usedNodePorts:        usedNodePorts,
		virtualConstants:     constantsToInsert,
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
