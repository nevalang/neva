package helper

import ts "github.com/emil14/neva/pkg/types"

func Def(body ts.Expr, params ...ts.Param) ts.Def {
	return ts.Def{
		Params: params,
		Body:   body,
	}
}

// Do not pass empty string as a name to avoid Body.Empty() == true
func DefWithoutBody(params ...ts.Param) ts.Def {
	return ts.Def{
		Params: params,
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

func Base(tt ...string) map[string]struct{} {
	m := make(map[string]struct{}, len(tt))
	for _, t := range tt {
		m[t] = struct{}{}
	}
	return m
}

func Param(name string, constr ts.Expr) ts.Param {
	return ts.Param{
		Name:       name,
		Constraint: constr,
	}
}

func ParamWithoutConstr(name string) ts.Param {
	return Param(name, ts.Expr{})
}
