package desugarer

import (
	"fmt"
	"sync/atomic"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

type handleStructSelectorsResult struct {
	connToReplace     src.Connection
	connToInsert      src.Connection
	constToInsertName string
	constToInsert     src.ConstDef
	nodeToInsert      src.Node
	nodeToInsertName  string
}

var selectorNodeRef = core.EntityRef{
	Pkg:  "builtin",
	Name: "Field",
}

var virtualSelectorsCount atomic.Uint64

// desugarStructSelectors replaces one connection with 2 connections and a node with const
func (d Desugarer) desugarStructSelectors(
	normConn src.NormalConnection,
) (handleStructSelectorsResult, *compiler.Error) {
	senderSide := normConn.SenderSide

	constCounter := virtualConstCount.Load()
	virtualConstCount.Store(constCounter + 1)
	constName := fmt.Sprintf("__const__%d", constCounter)

	counter := virtualSelectorsCount.Load()
	virtualSelectorsCount.Store(counter + 1)
	nodeName := fmt.Sprintf("__field__%d", counter)

	selectorNode := src.Node{
		Directives: map[src.Directive][]string{
			compiler.BindDirective: {constName},
		},
		EntityRef: selectorNodeRef,
	}

	// original connection must be replaced with two new connections, this is the first one
	connToReplace := src.Connection{
		Normal: &src.NormalConnection{
			SenderSide: src.ConnectionSenderSide{
				// preserve original sender
				PortAddr: senderSide.PortAddr,
				Const:    senderSide.Const,
				// but remove selectors in desugared version
				Selectors: nil,
			},
			ReceiverSide: src.ConnectionReceiverSide{
				Receivers: []src.ConnectionReceiver{
					{
						PortAddr: src.PortAddr{
							Node: nodeName, // point it to created selector node
							Port: "msg",
						},
					},
				},
				// don't forget sometimes we have both struct selectors and deferred connections
				DeferredConnections: normConn.ReceiverSide.DeferredConnections,
			},
		},
	}

	// and this is the second
	connToInsert := src.Connection{
		Normal: &src.NormalConnection{
			SenderSide: src.ConnectionSenderSide{
				PortAddr: &src.PortAddr{
					Node: nodeName, // created node received data from original sender and is now sending it further
					Port: "msg",
				},
				Selectors: nil, // no selectors in desugared version
			},
			ReceiverSide: normConn.ReceiverSide, // preserve original receivers
		},
	}

	constWithCfgMsg := d.createConstWithCfgMsgForSelectorNode(senderSide)

	return handleStructSelectorsResult{
		connToReplace:     connToReplace,
		connToInsert:      connToInsert,
		constToInsertName: constName,
		constToInsert:     constWithCfgMsg,
		nodeToInsertName:  nodeName,
		nodeToInsert:      selectorNode,
	}, nil
}

var (
	pathConstTypeExpr = ts.Expr{
		Inst: &ts.InstExpr{
			Ref: core.EntityRef{Pkg: "builtin", Name: "list"},
			Args: []ts.Expr{
				{
					Inst: &ts.InstExpr{
						Ref: core.EntityRef{Pkg: "builtin", Name: "string"},
					},
				},
			},
		},
	}
)

func (Desugarer) createConstWithCfgMsgForSelectorNode(senderSide src.ConnectionSenderSide) src.ConstDef {
	result := make([]src.ConstValue, 0, len(senderSide.Selectors))

	for _, selector := range senderSide.Selectors {
		result = append(result, src.ConstValue{
			Message: &src.MsgLiteral{
				Str: compiler.Pointer(selector),
			},
		})
	}

	return src.ConstDef{
		TypeExpr: pathConstTypeExpr,
		Value: src.ConstValue{
			Message: &src.MsgLiteral{
				List: result,
			},
		},
	}
}
