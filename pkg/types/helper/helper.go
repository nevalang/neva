package helper

import ts "github.com/emil14/neva/pkg/types"

func BaseDefWithRecursion(params ...ts.Param) ts.Def {
	return ts.Def{
		Params:           params,
		RecursionAllowed: true,
	}
}

// Do not pass empty string as a name to avoid Body.Empty() == true
func BaseDef(params ...ts.Param) ts.Def {
	return Def(ts.Expr{}, params...)
}

func Def(body ts.Expr, params ...ts.Param) ts.Def {
	return ts.Def{
		Params: params,
		Body:   body,
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
		Lit: ts.LitExpr{Enum: els},
	}
}

func Arr(size int, typ ts.Expr) ts.Expr {
	return ts.Expr{
		Lit: ts.LitExpr{
			Arr: &ts.ArrLit{Expr: typ, Size: size},
		},
	}
}

func Union(els ...ts.Expr) ts.Expr {
	if els == nil { // for !lit.Empty()
		els = []ts.Expr{}
	}
	return ts.Expr{
		Lit: ts.LitExpr{Union: els},
	}
}

func Rec(v map[string]ts.Expr) ts.Expr {
	if v == nil { // for !lit.Empty()
		v = map[string]ts.Expr{}
	}
	return ts.Expr{
		Lit: ts.LitExpr{
			Rec: v,
		},
	}
}

// Base without recursion. If you need recursion create map by hand.
func Base(tt ...string) map[string]bool {
	m := make(map[string]bool, len(tt))
	for _, t := range tt {
		m[t] = false
	}
	return m
}

func ParamWithoutConstr(name string) ts.Param {
	return Param(name, ts.Expr{})
}

func Param(name string, constr ts.Expr) ts.Param {
	return ts.Param{
		Name:  name,
		Const: constr,
	}
}

func Trace(ss ...string) ts.Trace {
	var t *ts.Trace
	for _, s := range ss {
		tmp := ts.NewTrace(t, s)
		t = &tmp
	}
	return *t
}
