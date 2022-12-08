package src

func NewLocalTypeRef(name string) TypeRef {
	return NewTypeRef(name, "")
}

func NewTypeRef(name, pkg string) TypeRef {
	return TypeRef{
		Pkg:  pkg,
		Name: name,
	}
}

func NewRefExpr(ref TypeRef) TypeExpr {
	return TypeExpr{Ref: ref}
}

func NewPort(expr TypeExpr, isArray bool) Port {
	return Port{
		TypeExpr: expr,
		IsArray:  isArray,
	}
}
