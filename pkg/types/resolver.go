package types

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/pkg/tools"
)

type Resolver struct {
	validator expressionValidator
	checker   subtypeChecker
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
func (r Resolver) Resolve( //nolint:funlen,gocognit // https://github.com/emil14/neva/issues/181
	expr Expr,
	scope map[string]Def,
	base map[string]bool, // true means recursion allowed
	trace Trace,
) (Expr, error) {
	if err := r.validator.Validate(expr); err != nil { // todo remove embedding
		return Expr{}, fmt.Errorf("%w: %v", ErrInvalidExpr, err)
	}

	switch expr.Lit.Type() {
	case EnumLitType:
		return expr, nil // nothing to resolve in enum
	case ArrLitType:
		resolvedArrType, err := r.Resolve(expr.Lit.Arr.Expr, scope, base, trace)
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrArrType, err)
		}
		return Expr{
			Lit: LitExpr{
				Arr: &ArrLit{resolvedArrType, expr.Lit.Arr.Size},
			},
		}, nil
	case UnionLitType:
		resolvedUnion := make([]Expr, 0, len(expr.Lit.Union))
		for _, unionEl := range expr.Lit.Union {
			resolvedEl, err := r.Resolve(unionEl, scope, base, trace)
			if err != nil {
				return Expr{}, fmt.Errorf("%w: %v", ErrUnionUnresolvedEl, err)
			}
			resolvedUnion = append(resolvedUnion, resolvedEl)
		}
		return Expr{
			Lit: LitExpr{Union: resolvedUnion},
		}, nil
	case RecLitType:
		resolvedStruct := make(map[string]Expr, len(expr.Lit.Rec))
		for field, fieldExpr := range expr.Lit.Rec {
			resolvedFieldExpr, err := r.Resolve(fieldExpr, scope, base, trace)
			if err != nil {
				return Expr{}, fmt.Errorf("%w: %v", ErrRecFieldUnresolved, err)
			}
			resolvedStruct[field] = resolvedFieldExpr
		}
		return Expr{
			Lit: LitExpr{Rec: resolvedStruct},
		}, nil
	}

	newTrace := Trace{
		prev: &trace,
		v:    expr.Inst.Ref,
	}

	if newTrace.prev != nil && newTrace.v == newTrace.prev.v { // check
		return Expr{}, fmt.Errorf("%w: trace: %v", ErrDirectRecursion, newTrace)
	}

	t := &newTrace
	for t.prev != nil {
		t = t.prev
		if newTrace.v == t.v {
			return
		}
	}

	// if err := r.CheckTrace(newTrace, base); err != nil {

	// }

	def, ok := scope[expr.Inst.Ref] // check that ref type exist
	if !ok {
		return Expr{}, fmt.Errorf("%w: %v", ErrUndefinedRef, expr.Inst.Ref)
	}

	// because case with generics must be checked
	// if r.isDirectSelfRef(def.Body, expr.Inst.Ref) { // move to validator? not sure because of how tests written
	// 	return Expr{}, fmt.Errorf("%w: %v", ErrDirectRecursion, def)
	// }

	if len(def.Params) != len(expr.Inst.Args) { // args must not be > than params to avoid bad case with constraint
		return Expr{}, fmt.Errorf(
			"%w, want %d, got %d", ErrInstArgsLen, len(def.Params), len(expr.Inst.Args),
		)
	}

	resolvedArgs := make([]Expr, 0, len(def.Params))
	for i, param := range def.Params { // resolve args and constrs and check their compatibility
		resolvedArg, err := r.Resolve(expr.Inst.Args[i], scope, base, newTrace)
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrUnresolvedArg, err)
		}
		resolvedArgs = append(resolvedArgs, resolvedArg)
		if param.Constraint.Empty() {
			continue
		}
		resolvedConstr, err := r.Resolve(param.Constraint, scope, base, newTrace)
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrConstr, err)
		}
		if err := r.checker.Check(resolvedArg, resolvedConstr); err != nil {
			return Expr{}, fmt.Errorf(" %w: %v", ErrIncompatArg, err)
		}
	}

	if _, ok := base[expr.Inst.Ref]; ok {
		return Expr{
			Inst: InstExpr{
				Ref:  expr.Inst.Ref,
				Args: resolvedArgs,
			},
		}, nil
	}

	newExpr := Expr{ // <- substitution
		Inst: InstExpr{
			Ref:  def.Body.Inst.Ref,
			Args: resolvedArgs,
		},
	}

	return r.Resolve(newExpr, scope, base, newTrace)
}

func (r Resolver) isDirectSelfRef(defBody Expr, exprRef string) bool {
	return !defBody.Inst.Empty() && defBody.Lit.Empty() && defBody.Inst.Ref == exprRef
}

// Trace is a linked-list for tracing resolving path.
type Trace struct {
	prev *Trace // prev == nil for first element
	v    string
}

func (t Trace) String() string {
	s := "[" + t.v
	for t.prev != nil {
		t = *t.prev
		s += ", " + t.v
	}
	return s + "]"
}

func NewDefaultResolver() Resolver {
	return Resolver{
		validator: Validator{},
		checker:   SubtypeChecker{},
	}
}

func MustNewResolver(v expressionValidator, c subtypeChecker) Resolver {
	tools.NilPanic(v, c)
	return Resolver{v, c}
}
