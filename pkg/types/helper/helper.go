package helper

import ts "github.com/emil14/neva/pkg/types"

// Do not pass empty string as a name to avoid Body.Empty() == true
func NativeDef(name string, params ...ts.Param) ts.Def {
	return ts.Def{
		Params: params,
		Body: InstExpr(name),
	}
}

// Do not pass empty string as a name to avoid inst.Empty() == true
func InstExpr(ref string, args ...ts.Expr) ts.Expr {
	return ts.Expr{
		Inst: ts.InstExpr{
			Ref:  ref,
			Args: args,
		},
	}
}

func EnumLitExpr(els ...string) ts.Expr {
	if els == nil { // for !lit.Empty()
		els = []string{}
	}
	return ts.Expr{
		Lit: ts.LiteralExpr{EnumLit: els},
	}
}

func UnionLitExpr(els ...ts.Expr) ts.Expr {
	if els == nil { // for !lit.Empty()
		els = []ts.Expr{}
	}
	return ts.Expr{
		Lit: ts.LiteralExpr{UnionLit: els},
	}
}

func RecLitExpr(v map[string]ts.Expr) ts.Expr {
	if v == nil { // for !lit.Empty()
		v = map[string]ts.Expr{}
	}
	return ts.Expr{
		Lit: ts.LiteralExpr{
			RecLit: v,
		},
	}
}
