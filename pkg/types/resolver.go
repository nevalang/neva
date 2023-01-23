// todo scope ins't validated?
package types

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/pkg/tools"
	"golang.org/x/exp/maps"
)

type Resolver struct {
	expressionValidator
	subtypeChecker
}

//go:generate mockgen -source $GOFILE -destination mocks_test.go -package ${GOPACKAGE}_test
type (
	expressionValidator interface {
		Validate(Expr) error // returns error if expression's invariant broken
	}
	subtypeChecker interface {
		Check(Expr, Expr) error // Returns error if first expression is not a subtype of second
	}
)

var (
	ErrInvalidExpr        = errors.New("expression must be valid in order to be resolved")
	ErrUndefinedRef       = errors.New("expression refers to type that is not presented in the scope")
	ErrInstArgsLen        = errors.New("inst cannot have more arguments than reference type has parameters")
	ErrIncompatArg        = errors.New("argument is not subtype of the parameter's contraint")
	ErrUnresolvedArg      = errors.New("can't resolve argument")
	ErrConstr             = errors.New("can't resolve constraint")
	ErrArrType            = errors.New("could not resolve array type")
	ErrUnionUnresolvedEl  = errors.New("can't resolve union element")
	ErrRecFieldUnresolved = errors.New("can't resolve record field")
	ErrDirectRecursion    = errors.New("type definition's body must not be directly self referenced to itself")
	ErrIndirectRecursion  = errors.New("type definition's body must not be indirectly self referenced to itself")
)

// Transforms one expression into another where all references points to native types.
// It's a recursive process where each step starts with validation. Invalid expression always leads to error.
// For inst expr it checks compatibility between args and params and returns error if some constraint isn't satisfied.
// Then it updates scope by adding params of ref type with resolved args as values to allow substitution later.
// Then it checks whether base type of current ref type is native type to terminate with nil err and resolved expr.
// For non-native types process starts from the beginning with updated scope. New scope will contain values for params.
// For lit exprs logic is the following: for enum do nothing (it's valid and not composite, there's nothing to resolve),
// for array resolve it's type, for record and union apply recursion for it's every field/element.
func (r Resolver) Resolve( //nolint:funlen
	expr Expr,
	scope map[string]Def,
	base map[string]struct{},
	// visited map[string]struct{},
) (Expr, error) {
	if err := r.Validate(expr); err != nil { // todo remove embedding
		return Expr{}, fmt.Errorf("%w: %v", ErrInvalidExpr, err)
	}

	switch expr.Lit.Type() { // resolve literal
	case EnumLitType:
		return expr, nil // nothing to resolve in enum
	case ArrLitType:
		resolvedArrType, err := r.Resolve(expr.Lit.Arr.Expr, scope, base)
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrArrType, err)
		}
		return ArrExpr(expr.Lit.Arr.Size, resolvedArrType), nil
	case UnionLitType:
		resolvedUnion := make([]Expr, 0, len(expr.Lit.Union))
		for _, unionEl := range expr.Lit.Union {
			resolvedEl, err := r.Resolve(unionEl, scope, base)
			if err != nil {
				return Expr{}, fmt.Errorf("%w: %v", ErrUnionUnresolvedEl, err)
			}
			resolvedUnion = append(resolvedUnion, resolvedEl)
		}
		return Union(resolvedUnion...), nil
	case RecLitType:
		resolvedStruct := make(map[string]Expr, len(expr.Lit.Rec))
		for field, fieldExpr := range expr.Lit.Rec {
			resolvedFieldExpr, err := r.Resolve(fieldExpr, scope, base)
			if err != nil {
				return Expr{}, fmt.Errorf("%w: %v", ErrRecFieldUnresolved, err)
			}
			resolvedStruct[field] = resolvedFieldExpr
		}
		return Rec(resolvedStruct), nil
	} // at this point we know it's an instantiation, not literal

	def, ok := scope[expr.Inst.Ref] // check that reference type exists
	if !ok {
		return Expr{}, fmt.Errorf("%w: %v", ErrUndefinedRef, expr.Inst.Ref)
	}

	isDefBodyInst := !def.Body.Inst.Empty() && def.Body.Lit.Empty() // check both because def body wasn't validated yet
	if isDefBodyInst && def.Body.Inst.Ref == expr.Inst.Ref {        // restrict unresolvable cases like t=t
		return Expr{}, fmt.Errorf("%w: %v", ErrDirectRecursion, def)
	}

	// check that args for every param is present
	if len(def.Params) != len(expr.Inst.Args) { // args must not be > than params to avoid bad case with constraint
		return Expr{}, fmt.Errorf(
			"%w, want %d, got %d", ErrInstArgsLen, len(def.Params), len(expr.Inst.Args),
		)
	}

	newScope := make(map[string]Def, len(scope)+len(def.Params)) // new scope will contain resolved args (shadow)
	maps.Copy(newScope, scope)
	resolvedArgs := make([]Expr, 0, len(def.Params)) // keep track of ordered resolved args if expr refers to native type

	for i, param := range def.Params { // resolve arguments and parameter's constraints to compare them
		resolvedArg, err := r.Resolve(expr.Inst.Args[i], scope, base)
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrUnresolvedArg, err)
		}

		resolvedArgs = append(resolvedArgs, resolvedArg)
		newScope[param.Name] = Def{Body: resolvedArg} // no params for types from args (substutution happens here)

		if param.Constraint.Empty() {
			continue
		}

		resolvedConstraint, err := r.Resolve(param.Constraint, scope, base) // should we resolve it here?
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrConstr, err)
		}
		if err := r.Check(resolvedArg, resolvedConstraint); err != nil { // compatibility check
			return Expr{}, fmt.Errorf(" %w: %v", ErrIncompatArg, err)
		}
	} // at this point we have resolved args that are compatible with their parameters and a new scope

	if _, ok := base[expr.Inst.Ref]; ok {
		return Inst(expr.Inst.Ref, resolvedArgs...), nil
	}

	return r.Resolve(def.Body, newScope, base)
}

func NewDefaultResolver() Resolver {
	return Resolver{
		expressionValidator: Validator{},
		subtypeChecker:      SubtypeChecker{},
	}
}

func MustNewResolver(v expressionValidator, c subtypeChecker) Resolver {
	tools.NilPanic(v, c)
	return Resolver{v, c}
}
