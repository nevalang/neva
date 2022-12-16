package types

type Def struct {
	Params []Param // Body can refer to these parameters
	Body   Expr    // Expression that must be resolved
}

type Param struct {
	Name       string // Must be unique among other type's parameters
	Constraint Expr   // Expression that must be resolved supertype of corresponding argument
}

// Either Instantiation or literal
type Expr struct {
	Literal       *LiteralExpr // If nil, then instantiation
	Instantiation InstantiationExpr
}

type InstantiationExpr struct {
	Ref  string // Must be in the scope
	Args []Expr // Every ref's parameter must have subtype argument
}

// Only one field must be initialized
type LiteralExpr struct {
	ArrLir   *ArrLit
	RecLit   map[string]Expr
	EnumLit  []string
	UnionLit []Expr
}

type ArrLit struct {
	Expr Expr
	Size uint8
}
