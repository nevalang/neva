// Package ts implements static type-system with generics and structural subtyping.
package ts

import "fmt"

type Def struct { // TODO add validation (what validation?)
	Params             []Param // Body can refer to these parameters
	BodyExpr           *Expr   // Empty body here means base type (TODO maybe change needed)
	IsRecursionAllowed bool    // Type can be used for recursive definition. Only base type can have true here
}

func (def Def) String() string {
	var params string

	if len(def.Params) > 0 {
		params += "<"
		for i, param := range def.Params {
			params += param.Name
			if param.Constr == nil {
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
	Constr *Expr  // Expression that must be resolved supertype of corresponding argument
}

// Instantiation or literal. Lit or Inst must be not nil, but not both
type Expr struct {
	Lit  *LitExpr
	Inst *InstExpr
}

// String formats expression in a TS manner
func (expr *Expr) String() string {
	if expr == nil {
		return "empty"
	}

	var str string

	switch expr.Lit.Type() {
	case ArrLitType:
		return fmt.Sprintf(
			"[%d]%s",
			expr.Lit.Arr.Size, expr.Lit.Arr.Expr.String(),
		)
	case EnumLitType:
		str += "{"
		for i, el := range expr.Lit.Enum {
			str += " " + el
			if i == len(expr.Lit.Enum)-1 {
				str += " "
			} else {
				str += ","
			}
		}
		return str + "}"
	case RecLitType:
		str += "{"
		count := 0
		for fieldName, fieldExpr := range expr.Lit.Rec {
			str += fmt.Sprintf(" %v %v", fieldName, fieldExpr)
			if count < len(expr.Lit.Rec)-1 {
				str += ","
			} else {
				str += " "
			}
			count++
		}
		return str + "}"
	case UnionLitType:
		for i, el := range expr.Lit.Union {
			str += el.String()
			if i < len(expr.Lit.Union)-1 {
				str += " | "
			}
		}
		return str
	}

	if len(expr.Inst.Args) == 0 {
		return expr.Inst.Ref
	}

	str = expr.Inst.Ref + "<"
	for i, arg := range expr.Inst.Args {
		str += arg.String()
		if i < len(expr.Inst.Args)-1 {
			str += ", "
		}
	}
	str += ">"

	return str
}

// Instantiation expression
type InstExpr struct {
	Ref  string // Must be in the scope
	Args []Expr // Every ref's parameter must have subtype argument
}

// Literal expression. Only one field must be initialized
type LitExpr struct {
	Arr   *ArrLit
	Rec   map[string]Expr
	Enum  []string
	Union []Expr
}

func (lit *LitExpr) Empty() bool {
	return lit == nil ||
		lit.Arr == nil &&
			lit.Rec == nil &&
			lit.Enum == nil &&
			lit.Union == nil
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
