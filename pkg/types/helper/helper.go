package helper

import "github.com/emil14/neva/pkg/types"

func NativeType(name string, params ...types.Param) types.Def {
	return types.Def{
		Params: params,
		Body: types.Expr{
			Inst: types.InstExpr{Ref: name},
		},
	}
}

func InstExpr(ref string, args ...types.Expr) types.Expr {
	return types.Expr{
		Inst: types.InstExpr{
			Ref:  ref,
			Args: args,
		},
	}
}
