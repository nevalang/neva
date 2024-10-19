package desugarer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

type desugarStructSelectorsResult struct {
	connToReplace src.Connection
	insert        src.Connection
}

var selectorNodeRef = core.EntityRef{
	Pkg:  "builtin",
	Name: "Field",
}

var virtualSelectorsCount uint64

// desugarStructSelectors replaces one connection with 2 connections and a node with const
func (d Desugarer) desugarStructSelectors(
	sender src.ConnectionSender,
	normConn src.NormalConnection,
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) (
	desugarStructSelectorsResult,
	*compiler.Error,
) {
	virtualConstCount++
	constName := fmt.Sprintf("__const__%d", virtualConstCount)

	virtualSelectorsCount++
	selectorNodeName := fmt.Sprintf("__field__%d", virtualSelectorsCount)

	selectorNode := src.Node{
		Directives: map[src.Directive][]string{
			compiler.BindDirective: {constName},
		},
		EntityRef: selectorNodeRef,
	}

	sender.Selectors = nil

	// original connection must be replaced with two new connections, this is the first one
	replace := src.Connection{
		Normal: &src.NormalConnection{
			SenderSide: []src.ConnectionSender{sender},
			ReceiverSide: []src.ConnectionReceiver{
				{
					PortAddr: &src.PortAddr{
						Node: selectorNodeName, // point it to created selector node
						Port: "data",
					},
				},
			},
		},
	}

	// and this is the second
	connToInsert := src.Connection{
		Normal: &src.NormalConnection{
			SenderSide: []src.ConnectionSender{
				{
					PortAddr: &src.PortAddr{
						Node: selectorNodeName, // created node received data from original sender and is now sending it further
						Port: "res",
					},
				},
			},
			ReceiverSide: normConn.ReceiverSide, // preserve original receivers
		},
	}

	nodesToInsert[selectorNodeName] = selectorNode
	constsToInsert[constName] = d.createSelectorCfgMsg(sender)

	return desugarStructSelectorsResult{
		connToReplace: replace,
		insert:        connToInsert,
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

func (Desugarer) createSelectorCfgMsg(senderSide src.ConnectionSender) src.Const {
	result := make([]src.ConstValue, 0, len(senderSide.Selectors))

	for _, selector := range senderSide.Selectors {
		result = append(result, src.ConstValue{
			Message: &src.MsgLiteral{
				Str: compiler.Pointer(selector),
			},
		})
	}

	return src.Const{
		TypeExpr: pathConstTypeExpr,
		Value: src.ConstValue{
			Message: &src.MsgLiteral{
				List: result,
			},
		},
	}
}
