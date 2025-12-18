// Package typesystem implements type-system with generics and structural subtyping.
// For convenience these structures have json tags (like `src` package).
// This is not clean architecture but it's very handy for LSP.
package typesystem

import (
	"sort"

	"github.com/nevalang/neva/internal/compiler/ast/core"
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

	if def.BodyExpr == nil {
		return params
	}
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
	case UnionLitType:
		// todo: keep deterministic order by sorting; ideally match source order which would require moving from map to slice (same for struct)
		count := 0
		str += "union {"
		// collect and sort tags for stable output
		var tags []string
		for tag := range expr.Lit.Union {
			tags = append(tags, tag)
		}
		sort.Strings(tags)
		for _, tag := range tags {
			tagExpr := expr.Lit.Union[tag]
			if count == 0 {
				str += " "
			}
			if tagExpr != nil {
				// check if this is a tag-only union (tag name matches type name)
				if tagExpr.Inst != nil && tagExpr.Inst.Ref.Name == tag {
					str += tag
				} else {
					str += tag + " " + tagExpr.String()
				}
			} else {
				str += tag
			}
			if count < len(tags)-1 {
				str += ", "
			}
			count++
		}
		if count > 0 {
			str += " "
		}
		str += "}"
		return str
	case StructLitType:
		// todo: keep deterministic order by sorting; ideally match source order which would require moving from map to slice
		str += "{"
		count := 0
		// collect and sort field names for stable output
		var fields []string
		for fieldName := range expr.Lit.Struct {
			fields = append(fields, fieldName)
		}
		sort.Strings(fields)
		for _, fieldName := range fields {
			fieldExpr := expr.Lit.Struct[fieldName]
			str += " " + fieldName + " " + fieldExpr.String()
			if count < len(fields)-1 {
				str += ","
			} else {
				str += " "
			}
			count++
		}
		return str + "}"
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
