package helper

import ts "github.com/emil14/neva/pkg/types"

// Do not pass empty string as a name to avoid Body.Empty() == true
func NativeDef(name string, params ...ts.Param) ts.Def {
	return ts.Def{
		Params: params,
		Body:   Inst(name),
	}
}

// Do not pass empty string as a name to avoid inst.Empty() == true
func Inst(ref string, args ...ts.Expr) ts.Expr {
	if args == nil {
		args = []ts.Expr{} // makes easier testing because resolver returns non-nil args for native types
	}
	return ts.Expr{
		Inst: ts.InstExpr{
			Ref:  ref,
			Args: args,
		},
	}
}

func Enum(els ...string) ts.Expr {
	if els == nil { // for !lit.Empty()
		els = []string{}
	}
	return ts.Expr{
		Lit: ts.LiteralExpr{Enum: els},
	}
}

func Arr(size int, typ ts.Expr) ts.Expr {
	return ts.Expr{
		Lit: ts.LiteralExpr{
			Arr: &ts.ArrLit{Expr: typ, Size: size},
		},
	}
}

func Union(els ...ts.Expr) ts.Expr {
	if els == nil { // for !lit.Empty()
		els = []ts.Expr{}
	}
	return ts.Expr{
		Lit: ts.LiteralExpr{Union: els},
	}
}

func Rec(v map[string]ts.Expr) ts.Expr {
	if v == nil { // for !lit.Empty()
		v = map[string]ts.Expr{}
	}
	return ts.Expr{
		Lit: ts.LiteralExpr{
			Rec: v,
		},
	}
}
