package src

func NewLocalPkgRef(name string) PkgRef {
	return NewPkgRef(name, "")
}

func NewLocalTypeRef(name string) TypeRef {
	return NewTypeRef(name, "")
}

func NewTypeRef(name, pkg string) TypeRef {
	return TypeRef{
		Pkg:  pkg,
		Name: name,
	}
}

func NewPkgRef(name, version string) PkgRef {
	return PkgRef{
		Name:    name,
		Version: version,
	}
}

func NewRefExpr(ref TypeRef) TypeExpr {
	return TypeExpr{Ref: ref}
}
