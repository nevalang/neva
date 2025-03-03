// Package typesystem implements type-system with generics and structural subtyping.
// For convenience these structures have json tags (like `src` package).
// This is not clean architecture but it's very handy for LSP.
package typesystem

import (
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

type Def struct {
	// Body can refer to these. Must be replaced with arguments while resolving
	Params []Param `json:"params,omitempty"`
	// Empty body means base type
	BodyExpr *Expr `json:"bodyExpr,omitempty"`
	// Meta can be used to store anything that can be useful for typesystem user. It is ignored by the typesystem itself.
	Meta core.Meta `json:"meta,omitempty"`
}

func (def Def) String() string {
	var params string

	params += "<"
	for i, param := range def.Params {
		params += param.Name
		params += " " + param.Constr.String()
		if i < len(def.Params)-1 {
			params += ", "
		}
	}
	params += ">"

	return params + " = " + def.BodyExpr.String()
}

type Param struct {
	Name   string `json:"name,omitempty"`   // Must be unique among other type's parameters
	Constr Expr   `json:"constr,omitempty"` // Expression that must be resolved supertype of corresponding argument
}

// Instantiation or literal. Lit or Inst must be not nil, but not both
type Expr struct {
	Lit  *LitExpr  `json:"lit,omitempty"`
	Inst *InstExpr `json:"inst,omitempty"`
	Meta core.Meta `json:"meta,omitempty"` // This field must be ignored by the typesystem and only used outside
}

// String formats expression in a TS manner
func (expr Expr) String() string {
	if expr.Inst == nil && expr.Lit == nil {
		return "empty"
	}

	var str string

	switch expr.Lit.Type() {
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
	case StructLitType:
		str += "{"
		count := 0
		for fieldName, fieldExpr := range expr.Lit.Struct {
			str += " " + fieldName + " " + fieldExpr.String()
			if count < len(expr.Lit.Struct)-1 {
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
		return expr.Inst.Ref.String()
	}

	str = expr.Inst.Ref.String()
	str += "<"

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
	Ref  core.EntityRef `json:"ref,omitempty"`  // Must be in the scope
	Args []Expr         `json:"args,omitempty"` // Every ref's parameter must have subtype argument
}

// Literal expression. Only one field must be initialized
type LitExpr struct {
	Struct map[string]Expr `json:"struct,omitempty"`
	Enum   []string        `json:"enum,omitempty"`
	Union  []Expr          `json:"union,omitempty"`
}

func (lit *LitExpr) Empty() bool {
	return lit == nil ||
		lit.Struct == nil &&
			lit.Enum == nil &&
			lit.Union == nil
}

// Always call Validate before
func (lit *LitExpr) Type() LiteralType {
	switch {
	case lit == nil:
		return EmptyLitType
	case lit.Struct != nil:
		return StructLitType
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
	StructLitType
	EnumLitType
	UnionLitType
)
