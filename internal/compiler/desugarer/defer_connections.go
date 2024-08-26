package desugarer

import (
	"fmt"
	"maps"
	"sync/atomic"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	"github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

var virtualBlockerNode = src.Node{
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

type desugarDeferredConnectionsResult struct {
	connsToInsert     []src.Connection
	receiversToInsert []src.ConnectionReceiver
	constsToInsert    map[string]src.Const
	nodesToInsert     map[string]src.Node
	nodesPortsUsed    nodePortsMap // (probably?) to generate "Del" instances where needed
}

var virtualBlockersCounter atomic.Uint64

func (d Desugarer) desugarDeferredConnections(
	originalConn src.NormalConnection,
	nodes map[string]src.Node,
	scope src.Scope,
) (desugarDeferredConnectionsResult, *compiler.Error) {
	deferredConnections := originalConn.ReceiverSide.DeferredConnections

	// recursively desugar every deferred connections
	handleNetResult, err := d.handleNetwork(
		deferredConnections,
		nodes,
		scope,
	)
	if err != nil {
		return desugarDeferredConnectionsResult{}, nil
	}

	// we want to return nodes created in recursive calls
	// as well as the onces created by us in this call
	nodesToInsert := maps.Clone(handleNetResult.nodesToInsert)

	// we going to replace all desugared deferreded connections with set of our connections
	connsToInsert := make([]src.Connection, 0, len(handleNetResult.desugaredConnections))

	// for every deferred connection we must do 4 things
	// 1) create virtual "blocker" node
	// 2) create connection from original sender to blocker:sig
	// 3) create connection from deferred sender to blocker:data
	// 4) create connection from blocker:data to every receiver in deferred connection

	// we gonna collect receivers for first connection instead of
	// creating several separate connections because that won't work
	receiversForOriginalSender := make([]src.ConnectionReceiver, 0, len(handleNetResult.desugaredConnections))

	for _, desugaredThenConn := range handleNetResult.desugaredConnections {
		deferredConnection := desugaredThenConn.Normal

		// 1) create and add virtual blocker node
		counter := virtualBlockersCounter.Load()
		virtualBlockersCounter.Store(counter + 1)
		virtualBlockerName := fmt.Sprintf("__lock__%d", counter)
		nodesToInsert[virtualBlockerName] = virtualBlockerNode

		// 2) create connection from original sender to blocker:sig
		receiversForOriginalSender = append(
			receiversForOriginalSender,
			src.ConnectionReceiver{
				PortAddr: src.PortAddr{
					Node: virtualBlockerName,
					Port: "sig",
				},
			},
		)

		connsToInsert = append(connsToInsert,
			// 3) create connection from deferred sender to blocker:data
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
			// 4) create connection from blocker:data to every receiver in deferred connection
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

	return desugarDeferredConnectionsResult{
		nodesToInsert:     nodesToInsert,
		connsToInsert:     connsToInsert,
		constsToInsert:    handleNetResult.constsToInsert,
		nodesPortsUsed:    handleNetResult.nodesPortsUsed,
		receiversToInsert: receiversForOriginalSender,
	}, nil
}
