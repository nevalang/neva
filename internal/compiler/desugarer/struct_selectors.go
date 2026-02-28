package desugarer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	ts "github.com/nevalang/neva/internal/compiler/typesystem"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
)

type desugarStructSelectorsResult struct {
	replace src.Connection
}

// desugarStructSelectors doesn't generate incoming connections for field node,
// it's responsibility of desugarChainConnection.
func (d *Desugarer) desugarStructSelectors(
	conn src.Connection, // sender here is selector (this is chained connection)
	nodesToInsert map[string]src.Node,
	constsToInsert map[string]src.Const,
) desugarStructSelectorsResult {
	locOnlyMeta := core.Meta{
		Location: conn.Senders[0].Meta.Location, // FIXME for some reason connection meta sometimes empty
	}

	d.virtualConstCount++
	constName := fmt.Sprintf("__const__%d", d.virtualConstCount)

	d.virtualSelectorsCount++
	selectorNodeName := fmt.Sprintf("__field__%d", d.virtualSelectorsCount)

	selectorNode := src.Node{
		Directives: map[src.Directive]string{compiler.BindDirective: constName},
		EntityRef: core.EntityRef{
			Pkg:  "builtin",
			Name: "Field",
			Meta: locOnlyMeta,
		},
		Meta: locOnlyMeta,
	}

	// struct selectors are discarded from this point
	replace := src.Connection{
		// created node will receive data from prev chain link
		Senders: []src.ConnectionSender{
			{
				PortAddr: &src.PortAddr{
					Node: selectorNodeName,
					Port: "res",
					Meta: locOnlyMeta,
				},
				Meta: locOnlyMeta,
			},
		},
		// and send it to original receiver side
		Receivers: conn.Receivers,
		Meta:      locOnlyMeta,
	}

	nodesToInsert[selectorNodeName] = selectorNode
	constsToInsert[constName] = d.createSelectorCfgMsg(conn.Senders[0])

	return desugarStructSelectorsResult{
		replace: replace,
	}
}

func pathConstTypeExpr() ts.Expr {
	return ts.Expr{
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
}

func (Desugarer) createSelectorCfgMsg(senderSide src.ConnectionSender) src.Const {
	result := make([]src.ConstValue, 0, len(senderSide.StructSelector))
	locOnlyMeta := core.Meta{
		Location: senderSide.Meta.Location,
	}

	for _, selector := range senderSide.StructSelector {
		selectorValue := selector
		result = append(result, src.ConstValue{
			Message: &src.MsgLiteral{
				Str:  &selectorValue,
				Meta: locOnlyMeta,
			},
		})
	}

	return src.Const{
		TypeExpr: pathConstTypeExpr(),
		Value: src.ConstValue{
			Message: &src.MsgLiteral{
				List: result,
				Meta: locOnlyMeta,
			},
		},
		Meta: locOnlyMeta,
	}
}
