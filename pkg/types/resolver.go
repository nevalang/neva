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
	ErrInvalidExpr                  = errors.New("expression must be valid in order to be resolved")
	ErrUndefinedRef                 = errors.New("expression refers to type that is not presented in the scope and args")
	ErrInstArgsLen                  = errors.New("inst must have same number of arguments as def has parameters")
	ErrIncompatArg                  = errors.New("argument is not subtype of the parameter's contraint")
	ErrUnresolvedArg                = errors.New("can't resolve argument")
	ErrConstr                       = errors.New("can't resolve constraint")
	ErrArrType                      = errors.New("could not resolve array type")
	ErrUnionUnresolvedEl            = errors.New("can't resolve union element")
	ErrRecFieldUnresolved           = errors.New("can't resolve record field")
	ErrDirectRecursion              = errors.New("type definition's body must not be directly self referenced to itself")
	ErrIndirectRecursion            = errors.New("type definition's body must not be indirectly self referenced to itself")
	ErrNotBaseTypeSupportsRecursion = errors.New("only base type definitions can have support for recursion")
)

// Resolve turn one expression into another where all references points to native types.
// It's a recursive process where each step starts with validation. Invalid expression always leads to error.
// For inst expr it checks compatibility between args and params and returns error if some constraint isn't satisfied.
// Then it updates scope by adding params of ref type with resolved args as values to allow substitution later.
// Then it checks whether base type of current ref type is native type to terminate with nil err and resolved expr.
// For non-native types process starts from the beginning with updated scope. New scope will contain values for params.
// For lit exprs logic is the following: for enum do nothing (it's valid and not composite, there's nothing to resolve),
// for array resolve it's type, for record and union apply recursion for it's every field/element.
func (r Resolver) Resolve( //nolint:funlen // https://github.com/emil14/neva/issues/181
	expr Expr, // expression that must be resolved
	scope map[string]Def, // immutable map of type definitions
	frame map[string]Def, // parameters mapped onto arguments at previous step
	trace TraceChain, // linked list to handle recursive expressions
) (Expr, error) {
	if err := r.validator.Validate(expr); err != nil { // todo remove embedding
		return Expr{}, fmt.Errorf("%w: %v", ErrInvalidExpr, err)
	}

	switch expr.Lit.Type() {
	case EnumLitType:
		return expr, nil // nothing to resolve in enum
	case ArrLitType:
		resolvedArrType, err := r.Resolve(expr.Lit.Arr.Expr, scope, frame, trace)
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
			resolvedEl, err := r.Resolve(unionEl, scope, frame, trace)
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
			resolvedFieldExpr, err := r.Resolve(fieldExpr, scope, frame, trace)
			if err != nil {
				return Expr{}, fmt.Errorf("%w: %v", ErrRecFieldUnresolved, err)
			}
			resolvedStruct[field] = resolvedFieldExpr
		}
		return Expr{
			Lit: LitExpr{Rec: resolvedStruct},
		}, nil
	}

	def, err := r.getDef(expr.Inst.Ref, frame, scope)
	if err != nil {
		return Expr{}, err
	}

	if len(def.Params) != len(expr.Inst.Args) { // args must not be > than params to avoid bad case with constraint
		return Expr{}, fmt.Errorf(
			"%w, want %d, got %d", ErrInstArgsLen, len(def.Params), len(expr.Inst.Args),
		)
	}

	newTrace := TraceChain{
		prev: &trace,
		v:    expr.Inst.Ref,
	}

	ret, err := r.checkRecursion(newTrace, scope, def) // TODO what about args (for loop below)?
	if err != nil {
		return Expr{}, fmt.Errorf("%w", err)
	} else if ret {
		return expr, nil
	}

	newFrame := make(map[string]Def, len(def.Params))
	resolvedArgs := make([]Expr, 0, len(expr.Inst.Args))
	for i, param := range def.Params { // resolve args and constrs and check their compatibility
		resolvedArg, err := r.Resolve(expr.Inst.Args[i], scope, frame, newTrace)
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrUnresolvedArg, err)
		}
		newFrame[param.Name] = Def{Body: resolvedArg} // no params for generics
		resolvedArgs = append(resolvedArgs, resolvedArg)
		if param.Constraint.Empty() {
			continue
		}
		resolvedConstr, err := r.Resolve(param.Constraint, scope, newFrame, newTrace) //nolint:lll // we pass newFrame because constr can refer type param
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrConstr, err)
		}
		if err := r.checker.Check(resolvedArg, resolvedConstr); err != nil {
			return Expr{}, fmt.Errorf(" %w: %v", ErrIncompatArg, err)
		}
	}

	if def.Body.Empty() {
		return Expr{
			Inst: InstExpr{
				Ref:  expr.Inst.Ref,
				Args: resolvedArgs,
			},
		}, nil
	}

	// TODO investigate possibility to replace "flat" arguments with resolved args
	return r.Resolve(def.Body, scope, newFrame, newTrace)
}

// getDef checks for def in args, then in scope and returns err if expr refers no nothing.
func (Resolver) getDef(ref string, args, scope map[string]Def) (Def, error) {
	def, ok := args[ref]
	if ok {
		return def, nil
	}

	def, ok = scope[ref]
	if !ok {
		return Def{}, fmt.Errorf("%w: %v", ErrUndefinedRef, ref)
	}

	return def, nil
}

// checkRecursion returns true and nil error for recursive expressions that should not go on next step of resolving.
// It returns false and nil err for non-recursive expressions with valid trace
// and false with non-nil err for bad recursion cases.
func (Resolver) checkRecursion(trace TraceChain, scope map[string]Def, def Def) (bool, error) {
	if !def.Body.Empty() && def.RecursionAllowed { // only base type can be used for recursion
		return false, fmt.Errorf("%w: %v", ErrNotBaseTypeSupportsRecursion, def)
	}

	if trace.prev != nil && trace.v == trace.prev.v {
		return true, fmt.Errorf("%w: trace: %v", ErrDirectRecursion, trace)
	}

	var isPrevAllowRecursion bool
	if trace.prev != nil {
		isPrevAllowRecursion = scope[trace.prev.v].RecursionAllowed
	}

	prev := trace.prev
	for prev != nil {
		if prev.v == trace.v && isPrevAllowRecursion {
			return true, nil
		}
		if prev.v == trace.v {
			return true, fmt.Errorf("%w: %v", ErrIndirectRecursion, trace)
		}
		prev = prev.prev
	}

	return false, nil
}

// TraceChain is a linked-list for tracing resolving path.
type TraceChain struct { // TODO make private
	prev *TraceChain // prev == nil for first element
	v    string
}

func (t TraceChain) String() string {
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
