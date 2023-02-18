package types

type Def struct { // TODO add validation
	Params             []Param // Body can refer to these parameters
	BodyExpr           Expr    // Empty body here means base type (TODO maybe change needed)
	IsRecursionAllowed bool    // Type can be used for recursive definition. Only base type can have this
}

func (def Def) String() string {
	var params string

	if len(def.Params) > 0 {
		params += "<"
		for i, param := range def.Params {
			params += param.Name
			if param.Constr.Empty() {
				continue
			}
			params += " " + param.Constr.String()
			if i < len(def.Params)-1 {
				params += ", "
			}
		}
		params += ">"
	}

	return params + " = " + def.BodyExpr.String()
}

type Param struct {
	Name   string // Must be unique among other type's parameters
	Constr Expr   // Expression that must be resolved supertype of corresponding argument
}

// Instantiation or literal
type Expr struct {
	Lit  LitExpr // If empty then expr is inst
	Inst InstExpr
}

// Empty returns true if inst and lit both empty
func (expr Expr) Empty() bool {
	return expr.Inst.Empty() && expr.Lit.Empty()
}

// Instantiation expression
type InstExpr struct {
	Ref  string // Must be in the scope
	Args []Expr // Every ref's parameter must have subtype argument
}

func (i InstExpr) Empty() bool {
	return i.Ref == "" && len(i.Args) == 0
}

// Literal expression. Only one field must be initialized
type LitExpr struct {
	Arr   *ArrLit
	Rec   map[string]Expr
	Enum  []string
	Union []Expr
}

// Helper to check that all lit's fields are nils. Doesn't care about validation
func (lit LitExpr) Empty() bool {
	return lit.Arr == nil && lit.Rec == nil && lit.Enum == nil && lit.Union == nil
}

// Always call Validate before
func (lit LitExpr) Type() LiteralType {
	switch {
	case lit.Arr != nil:
		return ArrLitType
	case lit.Rec != nil:
		return RecLitType
	case lit.Enum != nil:
		return EnumLitType
	case lit.Union != nil:
		return UnionLitType
	}
	return EmptyLitType // for inst or invalid lit
}

type LiteralType uint8

const (
	EmptyLitType LiteralType = iota
	ArrLitType
	RecLitType
	EnumLitType
	UnionLitType
)

type ArrLit struct {
	Expr Expr
	Size int
}
