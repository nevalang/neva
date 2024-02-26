package desugarer

import (
	"fmt"
	"maps"
	"sync/atomic"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
	"github.com/nevalang/neva/pkg/typesystem"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var virtualBlockerNode = src.Node{
	EntityRef: src.EntityRef{
		Pkg:  "builtin",
		Name: "Blocker",
	},
	TypeArgs: []typesystem.Expr{
		ts.Expr{
			Inst: &typesystem.InstExpr{
				Ref: src.EntityRef{Pkg: "builtin", Name: "any"},
			},
		},
	},
}

type handleThenConnectionsResult struct {
	desugaredConnections []src.Connection // these replaces original deferred connections
	virtualConstants     map[string]src.Const
	virtualNodes         map[string]src.Node
	usedNodesPorts       nodePortsMap
}

var virtualBlockersCounter atomic.Uint64

func (d Desugarer) handleDeferredConnections( //nolint:funlen
	sender src.ConnectionSenderSide, // who triggers deferred connections
	deferredConnections []src.Connection,
	nodes map[string]src.Node,
	scope src.Scope,
) (handleThenConnectionsResult, *compiler.Error) {
	// recursively desugar every deferred connections
	handleNetResult, err := d.handleNetwork(
		deferredConnections,
		nodes,
		scope,
	)
	if err != nil {
		return handleThenConnectionsResult{}, nil
	}

	// we want to return nodes created in recursive calls
	// as well as the onces created by us in this call
	virtualNodes := maps.Clone(handleNetResult.virtualNodes)

	// we going to replace all desugared deferreded connections with set of our connections
	virtualConns := make([]src.Connection, 0, len(handleNetResult.desugaredConnections))

	// for every deferred connection we must do 4 things
	// 1) create virtual "blocker" node
	// 2) create connection from original sender to blocker:sig
	// 3) create connection from deferred sender to blocker:data
	// 4) create connection from blocker:data to every receiver in deferred connection
	for _, desugaredThenConn := range handleNetResult.desugaredConnections {
		deferredConnection := desugaredThenConn.Normal

		// 1) create and add virtual blocker node
		counter := virtualBlockersCounter.Load()
		virtualBlockersCounter.Store(counter + 1)
		virtualBlockerName := fmt.Sprintf("virtual_blocker_%d", counter)
		virtualNodes[virtualBlockerName] = virtualBlockerNode

		// 2, 3 and 4 goes here
		virtualConns = append(virtualConns,
			// 2) original sender -> blocker:sig
			src.Connection{
				Normal: &src.NormalConnection{
					SenderSide: sender,
					ReceiverSide: src.ConnectionReceiverSide{
						Receivers: []src.ConnectionReceiver{
							{
								PortAddr: src.PortAddr{
									Node: virtualBlockerName,
									Port: "sig",
								},
							},
						},
					},
				},
			},
			// 2) deferred connection sender -> blocker:data
			src.Connection{
				Normal: &src.NormalConnection{
					SenderSide: deferredConnection.SenderSide,
					ReceiverSide: src.ConnectionReceiverSide{
						Receivers: []src.ConnectionReceiver{
							{
								PortAddr: src.PortAddr{
									Node: virtualBlockerName,
									Port: "data",
								},
							},
						},
					},
				},
			},
			// 3) blocker:data -> deferred connection receivers
			src.Connection{
				Normal: &src.NormalConnection{
					SenderSide: src.ConnectionSenderSide{
						PortAddr: &src.PortAddr{
							Node: virtualBlockerName,
							Port: "data",
						},
					},
					ReceiverSide: src.ConnectionReceiverSide{
						Receivers: deferredConnection.ReceiverSide.Receivers,
					},
				},
			},
		)
	}

	return handleThenConnectionsResult{
		virtualNodes:         virtualNodes,
		desugaredConnections: virtualConns,
		virtualConstants:     handleNetResult.virtualConstants,
		usedNodesPorts:       handleNetResult.usedNodePorts,
	}, nil
}
