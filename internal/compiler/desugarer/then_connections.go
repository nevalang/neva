package desugarer

import (
	"fmt"
	"maps"
	"slices"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
	"github.com/nevalang/neva/pkg/typesystem"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

var blockerNode = src.Node{
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

type handleThenConnsResult struct {
	extraConns     []src.Connection
	extraConsts    map[string]src.Const
	extraNodes     map[string]src.Node
	usedNodesPorts nodePortsMap
}

func (d Desugarer) handleThenConns( //nolint:funlen
	originalConn src.Connection,
	nodes map[string]src.Node,
	scope src.Scope,
) (handleThenConnsResult, *compiler.Error) {
	handleConnsResult, err := d.handleConns(originalConn.ReceiverSide.ThenConnections, nodes, scope)
	if err != nil {
		return handleThenConnsResult{}, nil
	}

	desugaredThenConns := handleConnsResult.desugaredConns
	extraNodes := maps.Clone(handleConnsResult.extraNodes)
	extraConns := slices.Clone(handleConnsResult.desugaredConns)

	for _, desugaredThenConn := range desugaredThenConns {
		blockerNodeName := fmt.Sprintf(
			"then_block_from_%v_to_%v_",
			originalConn.SenderSide.String(),
			desugaredThenConn.SenderSide.String(),
		)

		extraNodes[blockerNodeName] = blockerNode

		extraConns = append(
			extraConns,
			// original sender -> lock:sig
			src.Connection{
				SenderSide: originalConn.SenderSide,
				ReceiverSide: src.ConnectionReceiverSide{
					Receivers: []src.ConnectionReceiver{
						{
							PortAddr: src.PortAddr{
								Node: blockerNodeName,
								Port: "sig",
							},
						},
					},
				},
			},
			// then conn sender -> lock:data
			src.Connection{
				SenderSide: desugaredThenConn.SenderSide,
				ReceiverSide: src.ConnectionReceiverSide{
					Receivers: []src.ConnectionReceiver{
						{
							PortAddr: src.PortAddr{
								Node: blockerNodeName,
								Port: "data",
							},
						},
					},
				},
			},
			// lock:data -> { receivers... }
			src.Connection{
				SenderSide: src.ConnectionSenderSide{
					PortAddr: &src.PortAddr{
						Node: blockerNodeName,
						Port: "data",
					},
				},
				ReceiverSide: src.ConnectionReceiverSide{
					Receivers: desugaredThenConn.ReceiverSide.Receivers, // no nested then in desugared conn
				},
			},
		)
	}

	return handleThenConnsResult{
		extraNodes:     extraNodes,
		extraConns:     extraConns,
		extraConsts:    handleConnsResult.extraConsts,
		usedNodesPorts: handleConnsResult.usedNodePorts,
	}, nil
}
