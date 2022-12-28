package types

// Do not pass empty string as a name to avoid Body.Empty() == true
func NativeDef(name string, params ...Param) Def {
	return Def{
		Params: params,
		Body:   Inst(name),
	}
}

// Do not pass empty string as a name to avoid inst.Empty() == true
func Inst(ref string, args ...Expr) Expr {
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

func Enum(els ...string) Expr {
	if els == nil { // for !lit.Empty()
		els = []string{}
	}
	return Expr{
		Lit: LiteralExpr{Enum: els},
	}
}

func ArrExpr(size int, typ Expr) Expr {
	return Expr{
		Lit: LiteralExpr{
			Arr: &ArrLit{Expr: typ, Size: size},
		},
	}
}

func Union(els ...Expr) Expr {
	if els == nil { // for !lit.Empty()
		els = []Expr{}
	}
	return Expr{
		Lit: LiteralExpr{Union: els},
	}
}

func Rec(v map[string]Expr) Expr {
	if v == nil { // for !lit.Empty()
		v = map[string]Expr{}
	}
	return Expr{
		Lit: LiteralExpr{
			Rec: v,
		},
	}
}
