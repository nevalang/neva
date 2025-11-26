package desugarer

import (
	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/ast"
)

// handleConst handles case when constant has integer value and type is float.
func (d *Desugarer) handleConst(constant src.Const) (src.Const, error) {
	if constant.Value.Message == nil {
		return constant, nil
	}
	if constant.TypeExpr.String() != "float" {
		return constant, nil
	}
	if constant.Value.Message.Float != nil {
		return constant, nil
	}
	return src.Const{
		TypeExpr: constant.TypeExpr,
		Value: src.ConstValue{
			Message: &src.MsgLiteral{
				Float: compiler.Pointer(float64(*constant.Value.Message.Int)),
				Meta:  constant.Meta,
			},
		},
	}, nil
}
