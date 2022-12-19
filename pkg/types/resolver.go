package types

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/pkg/tools"
	"golang.org/x/exp/maps"
)

type Resolver struct {
	// table map[string]Def todo?
	validator
	checker
}

//go:generate mockgen -source $GOFILE -destination mocks_test.go -package ${GOPACKAGE}_test
type (
	validator interface {
		Validate(Expr) error // returns error if expression's invariant broken
	}
	checker interface {
		SubTypeCheck(Expr, Expr) error // Returns nil if first expr is subtype of the second
	}
)

var ErrInvalidExpr = errors.New("expr must be valid to be resolved")

// Transforms one expression into another where all references points to native types.
// It's a recursive process where each step starts with validation. Invalid expression always leads to error.
// For inst expr it checks compatibility between args and params and returns error if some constraint isn't satisfied.
// Then it updates scope by adding params of ref type with resolved args as values to allow substitution later.
// Then it checks whether base type of current ref type is native type to terminate with nil err and resolved expr.
// For non-native types process starts from the beginning with updated scope. New scope will contain values for params.
// For lit exprs logic is the following: for enum do nothing (it's valid and not composite, there's nothing to resolve),
// for array resolve it's type, for record and union apply recursion for it's every field/element.
func (r Resolver) Resolve(expr Expr, scope map[string]Def) (Expr, error) { //nolint:funlen
	if err := r.Validate(expr); err != nil {
		return Expr{}, fmt.Errorf("%w: %v", ErrInvalidExpr, err)
	}

	switch expr.Lit.Type() { // resolve literal
	case EnumLitType:
		return expr, nil // nothing to resolve in enum
	case ArrLitType:
		resolvedArrType, err := r.Resolve(expr.Lit.ArrLit.Expr, scope)
		if err != nil {
			return Expr{}, fmt.Errorf("invalid expr: %w", err)
		}
		return Expr{
			Lit: LiteralExpr{
				ArrLit: &ArrLit{resolvedArrType, expr.Lit.ArrLit.Size},
			},
		}, nil
	case UnionLitType:
		resolvedUnion := make([]Expr, 0, len(expr.Lit.UnionLit))
		for _, unionEl := range expr.Lit.UnionLit {
			resolvedEl, err := r.Resolve(unionEl, scope)
			if err != nil {
				return Expr{}, err
			}
			resolvedUnion = append(resolvedUnion, resolvedEl)
		}
		return Expr{
			Lit: LiteralExpr{UnionLit: resolvedUnion},
		}, nil
	case RecLitType:
		resolvedStruct := make(map[string]Expr, len(expr.Lit.RecLit))
		for field, fieldExpr := range expr.Lit.RecLit {
			resolvedFieldExpr, err := r.Resolve(fieldExpr, scope)
			if err != nil {
				return Expr{}, errors.New("")
			}
			resolvedStruct[field] = resolvedFieldExpr
		}
		return Expr{
			Lit: LiteralExpr{RecLit: resolvedStruct},
		}, nil
	}

	refType, ok := scope[expr.Inst.Ref] // check that reference type exists
	if !ok {
		return Expr{}, errors.New("ref type not found in scope")
	}

	if len(refType.Params) > len(expr.Inst.Args) { // check that generic args for every param is present
		return Expr{}, fmt.Errorf(
			"expr must have at least %d arguments, got %d",
			len(refType.Params), len(expr.Inst.Args),
		)
	}

	newScope := make(map[string]Def, len(scope)+len(refType.Params)) // new scope contains resolved params (shadow)
	maps.Copy(newScope, scope)
	resolvedArgs := make([]Expr, 0, len(refType.Params))

	for i, param := range refType.Params {
		resolvedArg, err := r.Resolve(expr.Inst.Args[i], scope)
		if err != nil {
			return Expr{}, errors.New("")
		}

		resolvedConstraint, err := r.Resolve(param.Constraint, scope) // should we resolve it here?
		if err != nil {
			return Expr{}, errors.New("")
		}

		if err := r.SubTypeCheck(resolvedArg, resolvedConstraint); err != nil { // compatibility check
			return Expr{}, fmt.Errorf("arg not subtype of constraint: %w", err)
		}

		resolvedArgs = append(resolvedArgs, resolvedArg)
		newScope[param.Name] = Def{
			Params: nil, // we don't refer generics with another generics inside
			Body:   resolvedArg,
		}
	}

	if refType.Body.Lit.Empty() { // reference type's body is an instantiation
		baseType, ok := scope[refType.Body.Inst.Ref]
		if !ok {
			return Expr{}, errors.New("")
		}
		if expr.Inst.Ref == baseType.Body.Inst.Ref { // direct self reference = native instantiation
			return Expr{
				Inst: InstExpr{
					Ref:  refType.Body.Inst.Ref,
					Args: resolvedArgs,
				},
			}, nil
		}
	}

	return r.Resolve(refType.Body, newScope) // it's not a native type and not literal - next step is needed
}

func DefaultResolver() Resolver {
	return Resolver{
		validator: nil,
		checker:   nil,
	}
}

// Allowes to pass custom validator and subtype checker
func NewResolver(v validator, c checker) Resolver {
	tools.PanicOnNil(v, c)
	return Resolver{v, c}
}
