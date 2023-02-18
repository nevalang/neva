package types

// Helper is just a namespace for helper functions to avoid conflicts with entity types
type Helper struct{}

func (h Helper) BaseDefWithRecursion(params ...Param) Def {
	return Def{
		Params:             params,
		IsRecursionAllowed: true,
	}
}

// Do not pass empty string as a name to avoid Body.Empty() == true
func (h Helper) BaseDef(params ...Param) Def {
	return Def{Params: params}
}

func (h Helper) Def(body Expr, params ...Param) Def {
	return Def{
		Params: params,
		Body:   body,
	}
}

// Do not pass empty string as a name to avoid inst.Empty() == true
func (h Helper) Inst(ref string, args ...Expr) Expr {
	if args == nil {
		args = []Expr{} // makes easier testing because resolver returns non-nil args for native types
	}
	return Expr{
		Inst: InstExpr{
			Ref:  ref,
			Args: args,
		},
	}
}

func (h Helper) Enum(els ...string) Expr {
	if els == nil { // for !lit.Empty()
		els = []string{}
	}
	return Expr{
		Lit: LitExpr{Enum: els},
	}
}

func (h Helper) Arr(size int, typ Expr) Expr {
	return Expr{
		Lit: LitExpr{
			Arr: &ArrLit{Expr: typ, Size: size},
		},
	}
}

func (h Helper) Union(els ...Expr) Expr {
	if els == nil { // for !lit.Empty()
		els = []Expr{}
	}
	return Expr{
		Lit: LitExpr{Union: els},
	}
}

func (h Helper) Rec(rec map[string]Expr) Expr {
	if rec == nil { // for !lit.Empty()
		rec = map[string]Expr{}
	}
	return Expr{
		Lit: LitExpr{
			Rec: rec,
		},
	}
}

func (h Helper) ParamWithNoConstr(name string) Param {
	return h.Param(name, Expr{})
}

func (h Helper) Param(name string, constr Expr) Param {
	return Param{
		Name:   name,
		Constr: constr,
	}
}

func (h Helper) Trace(ss ...string) Trace {
	var t *Trace
	for _, s := range ss {
		tmp := NewTrace(t, s)
		t = &tmp
	}
	return *t
}
