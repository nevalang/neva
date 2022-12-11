package types

import (
	"errors"

	"golang.org/x/exp/maps"
)

// Def with empty typeExpr and empty structFields is a native builtin Def
type Def struct { // l<t> = list<t> || l<t> = { foo t }
	Params []string // any type can have it
	Expr   Expr     // not for struct type (just "Type" in go spec)
}

// Expr is something that can be resolved and checked for compatibility with other expr.
// It's either Instantiation or literal. If it's a literal it's array, struct, enum or union
type Expr struct {
	Instantiation InstantiationExpr // (indirect recursion)
	ArrLir        *ArrLit
	RecLit        map[string]Expr
	EnumLit       []string
	UnionLit      []Expr
}

type ArrLit struct {
	TypeExpr Expr
	Size     uint8
}

// InstantiationExpr is an expression that always contains reference to some type and arguments if types has parameters.
type InstantiationExpr struct { // list<list<int>>
	Ref  string
	Args []Expr // can contain refs to generics, can be empty (indirect recursion!)
}

// resolve ...
func resolve(expr Expr, scope map[string]Def) (Expr, error) { // FIXME add support for recursive types (especially structs)
	if expr.EnumLit != nil {
		return expr, nil
	}

	if expr.UnionLit != nil {
		resolvedUnion := make([]Expr, 0, len(expr.UnionLit))
		for _, el := range expr.UnionLit {
			resolvedEl, err := resolve(el, scope)
			if err != nil {
				return Expr{}, err
			}
			resolvedUnion = append(resolvedUnion, resolvedEl)
		}
		return Expr{UnionLit: resolvedUnion}, nil
	}

	if expr.RecLit != nil {
		resolvedStruct := make(map[string]Expr, len(expr.RecLit))
		for field, expr := range expr.RecLit {
			resolvedFieldExpr, err := resolve(expr, scope)
			if err != nil {
				return Expr{}, errors.New("")
			}
			resolvedStruct[field] = resolvedFieldExpr
		}
		return Expr{
			RecLit: resolvedStruct,
		}, nil
	}

	refType, ok := scope[expr.Instantiation.Ref] // check that reference type exists
	if !ok {
		return Expr{}, errors.New("")
	}

	// check that generic args for every param is present
	if len(refType.Params) > len(expr.Instantiation.Args) { // compare equality? structural typing? linting?
		return Expr{}, errors.New("")
	}

	resolvedArgs := make([]Expr, 0, len(refType.Params))
	newScope := make(map[string]Def, len(scope)+len(refType.Params)) // new scope contains resolved params (shadow)
	// optimized for concurrency (is there better way?)
	maps.Copy(newScope, scope)
	for i, param := range refType.Params {
		resolvedArg, err := resolve(expr.Instantiation.Args[i], scope)
		if err != nil {
			return Expr{}, errors.New("")
		}
		resolvedArgs = append(resolvedArgs, resolvedArg)
		newScope[param] = Def{
			Params: nil, // we don't refer generics with another generics inside!
			Expr:   resolvedArg,
		}
	}

	if refType.Expr.RecLit == nil { // reference type's body is an application, not a struct definition
		baseType, ok := scope[refType.Expr.Instantiation.Ref] // FIXME not work structs
		if !ok {
			return Expr{}, errors.New("")
		}
		if expr.Instantiation.Ref == baseType.Expr.Instantiation.Ref {
			return Expr{
				Instantiation: InstantiationExpr{
					Ref:  refType.Expr.Instantiation.Ref,
					Args: resolvedArgs,
				},
				RecLit: nil, // todo
			}, nil
		}
	}

	return resolve(refType.Expr, newScope) // if it's not a native type and not a struct, then do recursive
}

func (expr Expr) String() string {
	var s string

	if expr.RecLit != nil {
		s += "{"
		for fieldName, fieldExpr := range expr.RecLit {
			s += " " + fieldName + ": " + fieldExpr.String() + " "
		}
		s += "}"
		return s
	}

	if len(expr.Instantiation.Args) == 0 {
		return expr.Instantiation.Ref
	}

	s = expr.Instantiation.Ref + "<"
	for i, arg := range expr.Instantiation.Args {
		s += arg.String()
		if i < len(expr.Instantiation.Args)-1 {
			s += ", "
		}
	}
	s += ">"

	return s
}
