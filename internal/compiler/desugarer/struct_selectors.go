package desugarer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

type desugarStructSelectorsResult struct {
	replace src.Connection
}

var selectorNodeRef = core.EntityRef{
	Pkg:  "builtin",
	Name: "Field",
}

var virtualSelectorsCount uint64

// desugarStructSelectors doesn't generate incoming connections for field node,
// it's responsibility of desugarChainConnection.
func (d Desugarer) desugarStructSelectors(
	normConn src.NormalConnection, // sender here is selector (this is chained connection)
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

	// struct selectors are discarded from this point
	replace := src.Connection{
		Normal: &src.NormalConnection{
			// created node will receive data from prev chain link
			SenderSide: []src.ConnectionSender{
				{
					PortAddr: &src.PortAddr{
						Node: selectorNodeName,
						Port: "res",
					},
				},
			},
			// and send it to original receiver side
			ReceiverSide: normConn.ReceiverSide,
		},
	}

	nodesToInsert[selectorNodeName] = selectorNode
	constsToInsert[constName] = d.createSelectorCfgMsg(normConn.SenderSide[0])

	return desugarStructSelectorsResult{
		replace: replace,
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
	result := make([]src.ConstValue, 0, len(senderSide.StructSelector))

	for _, selector := range senderSide.StructSelector {
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
