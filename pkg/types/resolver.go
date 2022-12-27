// todo scope ins't validated?
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
		SubtypeCheck(Expr, Expr) error // Returns nil if first expr is subtype of the second
	}
)

var (
	ErrInvalidExpr        = errors.New("expression must be valid in order to be resolved")
	ErrUndefinedRef       = errors.New("expression refers to type that is not presented in the scope")
	ErrInstArgsLen        = errors.New("inst cannot have more arguments than reference type has parameters")
	ErrIncompatArg        = errors.New("argument is not subtype of the parameter's contraint")
	ErrConstr             = errors.New("can't resolve constraint")
	ErrNoBaseType         = errors.New("definition's body refers to type that is not in the scope")
	ErrArrType            = errors.New("could not resolve array type")
	ErrUnionUnresolvedEl  = errors.New("can't resolve union element")
	ErrRecFieldUnresolved = errors.New("can't resolve record field")
)

// Transforms one expression into another where all references points to native types.
// It's a recursive process where each step starts with validation. Invalid expression always leads to error.
// For inst expr it checks compatibility between args and params and returns error if some constraint isn't satisfied.
// Then it updates scope by adding params of ref type with resolved args as values to allow substitution later.
// Then it checks whether base type of current ref type is native type to terminate with nil err and resolved expr.
// For non-native types process starts from the beginning with updated scope. New scope will contain values for params.
// For lit exprs logic is the following: for enum do nothing (it's valid and not composite, there's nothing to resolve),
// for array resolve it's type, for record and union apply recursion for it's every field/element.
func (r Resolver) Resolve(expr Expr, scope map[string]Def) (Expr, error) { //nolint:funlen,gocognit
	if err := r.Validate(expr); err != nil {
		return Expr{}, fmt.Errorf("%w: %v", ErrInvalidExpr, err)
	}

	switch expr.Lit.Type() { // resolve literal
	case EnumLitType:
		return expr, nil // nothing to resolve in enum
	case ArrLitType:
		resolvedArrType, err := r.Resolve(expr.Lit.ArrLit.Expr, scope)
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrArrType, err)
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
				return Expr{}, fmt.Errorf("%w: %v", ErrUnionUnresolvedEl, err)
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
				return Expr{}, fmt.Errorf("%w: %v", ErrRecFieldUnresolved, err)
			}
			resolvedStruct[field] = resolvedFieldExpr
		}
		return Expr{
			Lit: LiteralExpr{RecLit: resolvedStruct},
		}, nil
	}

	def, ok := scope[expr.Inst.Ref] // check that reference type exists
	if !ok {
		return Expr{}, fmt.Errorf("%w: %v", ErrUndefinedRef, expr.Inst.Ref)
	}

	// check that args for every param is present
	if len(def.Params) != len(expr.Inst.Args) { // args must not be > than params to avoid bad case with constraint
		return Expr{}, fmt.Errorf(
			"%w, want %d, got %d", ErrInstArgsLen, len(def.Params), len(expr.Inst.Args),
		)
	}

	newScope := make(map[string]Def, len(scope)+len(def.Params)) // new scope contains resolved params (shadow)
	maps.Copy(newScope, scope)
	resolvedArgs := make([]Expr, 0, len(def.Params)) // in case of native type

	for i, param := range def.Params {
		resolvedArg, err := r.Resolve(expr.Inst.Args[i], scope)
		if err != nil {
			return Expr{}, errors.New("")
		}

		if !param.Constraint.Empty() {
			resolvedConstraint, err := r.Resolve(param.Constraint, scope) // should we resolve it here?
			if err != nil {
				return Expr{}, fmt.Errorf("%w: %v", ErrConstr, err)
			}
			if err := r.SubtypeCheck(resolvedArg, resolvedConstraint); err != nil { // compatibility check
				return Expr{}, fmt.Errorf(" %w: %v", ErrIncompatArg, err)
			}
		}

		resolvedArgs = append(resolvedArgs, resolvedArg)
		newScope[param.Name] = Def{Body: resolvedArg} // no params for types from args
	}

	if def.Body.Lit.Empty() { // reference type's body is an instantiation
		baseType, ok := scope[def.Body.Inst.Ref]
		if !ok {
			return Expr{}, fmt.Errorf("%w: %v", ErrNoBaseType, def.Body.Inst.Ref)
		}
		if expr.Inst.Ref == baseType.Body.Inst.Ref { // direct self reference = native instantiation
			return Expr{
				Inst: InstExpr{
					Ref:  def.Body.Inst.Ref,
					Args: resolvedArgs,
				},
			}, nil
		}
	}

	return r.Resolve(def.Body, newScope) // it's not a native type and not literal - next step is needed
}

func NewDefaultResolver() Resolver {
	return Resolver{
		validator: Validator{},
		checker:   SubTypeChecker{},
	}
}

func MustNewResolver(v validator, c checker) Resolver {
	tools.PanicOnNil(v, c)
	return Resolver{v, c}
}
