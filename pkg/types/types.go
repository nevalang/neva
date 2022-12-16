package types

type Def struct {
	Params []Param // Body can refer to these parameters
	Body   Expr    // Expression that must be resolved
}

type Param struct {
	Name       string // Must be unique among other type's parameters
	Constraint Expr   // Expression that must be resolved supertype of corresponding argument
}

// Instantiation or literal
type Expr struct {
	Literal       LiteralExpr // If empty then instantiation
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

func (lit LiteralExpr) Empty() bool {
	return lit.ArrLir == nil && lit.RecLit == nil && lit.EnumLit == nil && lit.UnionLit == nil
}

func (lit LiteralExpr) Type() LiteralType {
	switch {
	case lit.ArrLir != nil:
		return ArrLitType
	case lit.RecLit != nil:
		return RecLitType
	case lit.EnumLit != nil:
		return EnumLitType
	case lit.UnionLit != nil:
		return UnionLitType
	}
	return UnknownLitType
}

type LiteralType uint8

const (
	UnknownLitType LiteralType = iota
	ArrLitType     LiteralType = iota
	RecLitType     LiteralType = iota
	EnumLitType    LiteralType = iota
	UnionLitType   LiteralType = iota
)

type ArrLit struct {
	Expr Expr
	Size uint8
}
