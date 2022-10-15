package src

type TypeExpr struct {
	Struct StructTypeExpr

	Ref     TypeRef
	RefArgs []TypeExpr
}

type TypeRef struct {
	Pkg  string
	Name string
}

type Type struct {
	Generics   []string
	Expr       TypeExpr
	StructExpr StructTypeExpr
}

type StructTypeExpr map[string]TypeExpr

func BuiltinTypes() map[string]Type { // move out?
	return map[string]Type{
		"bool":  {},
		"int":   {},
		"float": {},
		"str":   {},
		"list":  {Generics: []string{"t"}},
		"dict":  {Generics: []string{"t"}},
	}
}
