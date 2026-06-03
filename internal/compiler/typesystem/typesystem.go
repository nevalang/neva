// Package typesystem implements type-system with generics and structural subtyping.
// For convenience these structures have json tags (like `src` package).
// This is not clean architecture but it's very handy for LSP.
package typesystem

import (
	"sort"
	"strings"

	"github.com/nevalang/neva/pkg/core"
)

type Def struct {
	BodyExpr *Expr     `json:"bodyExpr,omitempty"`
	Params   []Param   `json:"params,omitempty"`
	Meta     core.Meta `json:"meta"`
}

func (def Def) String() string {
	var params strings.Builder

	params.WriteString("<")
	for i, param := range def.Params {
		params.WriteString(param.Name)
		params.WriteString(" " + param.Constr.String())
		if i < len(def.Params)-1 {
			params.WriteString(", ")
		}
	}
	params.WriteString(">")

	return params.String() + " = " + def.BodyExpr.String()
}

type Param struct {
	Name   string `json:"name,omitempty"` // Must be unique among other type's parameters
	Constr Expr   `json:"constr"`         // Expression that must be resolved supertype of corresponding argument
}

// Instantiation or literal. Lit or Inst must be not nil, but not both
type Expr struct {
	Lit  *LitExpr  `json:"lit,omitempty"`
	Inst *InstExpr `json:"inst,omitempty"`
	Meta core.Meta `json:"meta"` // This field must be ignored by the typesystem and only used outside
}

// String formats expression in a TS manner.
func (expr Expr) String() string {
	if expr.Inst == nil && expr.Lit == nil {
		return "empty"
	}

	if expr.Lit != nil {
		return (&expr).stringLit()
	}

	if len(expr.Inst.Args) == 0 {
		return expr.Inst.Ref.String()
	}

	return (&expr).stringInst()
}

func (expr *Expr) stringLit() string {
	switch expr.Lit.Type() {
	case EmptyLitType:
		return "empty"
	case UnionLitType:
		return expr.stringUnionLit()
	case StructLitType:
		return expr.stringStructLit()
	default:
		return "empty"
	}
}

func (expr *Expr) stringUnionLit() string {
	var str strings.Builder

	str.WriteString("union {")
	tags := sortedKeysFromUnion(expr.Lit.Union)
	for tagIndex, tag := range tags {
		if tagIndex == 0 {
			str.WriteString(" ")
		}

		tagExpr := expr.Lit.Union[tag]
		if tagExpr == nil || tagExpr.Inst != nil && tagExpr.Inst.Ref.Name == tag {
			str.WriteString(tag)
		} else {
			str.WriteString(tag)
			str.WriteString(" ")
			str.WriteString(tagExpr.String())
		}

		if tagIndex < len(tags)-1 {
			str.WriteString(", ")
		}
	}
	if len(tags) > 0 {
		str.WriteString(" ")
	}
	str.WriteString("}")

	return str.String()
}

func (expr *Expr) stringStructLit() string {
	var str strings.Builder

	str.WriteString("{")
	fields := sortedKeysFromStruct(expr.Lit.Struct)
	for i, fieldName := range fields {
		str.WriteString(" ")
		str.WriteString(fieldName)
		str.WriteString(" ")
		str.WriteString(expr.Lit.Struct[fieldName].String())
		if i < len(fields)-1 {
			str.WriteString(",")
		} else {
			str.WriteString(" ")
		}
	}
	str.WriteString("}")

	return str.String()
}

func (expr *Expr) stringInst() string {
	var str strings.Builder

	str.WriteString(expr.Inst.Ref.String())
	str.WriteString("<")
	for i, arg := range expr.Inst.Args {
		if i > 0 {
			str.WriteString(", ")
		}
		str.WriteString(arg.String())
	}
	str.WriteString(">")

	return str.String()
}

func sortedKeysFromUnion(union map[string]*Expr) []string {
	tags := make([]string, 0, len(union))
	for tag := range union {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	return tags
}

func sortedKeysFromStruct(structFields map[string]Expr) []string {
	fields := make([]string, 0, len(structFields))
	for fieldName := range structFields {
		fields = append(fields, fieldName)
	}
	sort.Strings(fields)

	return fields
}

// Instantiation expression
type InstExpr struct {
	Args []Expr         `json:"args,omitempty"`
	Ref  core.EntityRef `json:"ref"`
}

// Literal expression. Only one field must be initialized
type LitExpr struct {
	Struct map[string]Expr  `json:"struct,omitempty"`
	Union  map[string]*Expr `json:"union,omitempty"` // tag -> constraint
}

func (lit *LitExpr) Empty() bool {
	return lit == nil ||
		lit.Struct == nil &&
			lit.Union == nil
}

// Always call Validate before
func (lit *LitExpr) Type() LiteralType {
	switch {
	case lit == nil:
		return EmptyLitType
	case lit.Struct != nil:
		return StructLitType
	case lit.Union != nil:
		return UnionLitType
	}
	return EmptyLitType // for inst or invalid lit
}

type LiteralType uint8

const (
	EmptyLitType LiteralType = iota
	StructLitType
	UnionLitType
)
