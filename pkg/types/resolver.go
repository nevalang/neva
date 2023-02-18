package types

import (
	"errors"
	"fmt"

	"github.com/emil14/neva/pkg/tools"
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
	ErrNotBaseTypeSupportsRecursion = errors.New("only base type definitions can have support for recursion")
	ErrValidator                    = errors.New("validator implementation must not allow empty literals")
	ErrTerminator                   = errors.New("recursion terminator")
)

// Resolver transforms expression it into a form where all references it contains points to resolved expressions.
type Resolver struct {
	validator  exprValidator       // Check if expression invalid before resolving it
	comparator compatChecker       // Compare arguments with constraints
	terminator recursionTerminator // Don't stuck in a loop
}

//go:generate mockgen -source $GOFILE -destination mocks_test.go -package ${GOPACKAGE}_test
type (
	exprValidator interface {
		Validate(Expr) error // non-recursive validation of an expression
	}
	compatChecker interface {
		Check(Expr, Trace, Expr, Trace, map[string]Def) error // error should be nil if first expr is subtype of second
	}
	recursionTerminator interface {
		ShouldTerminate(Trace, map[string]Def) (bool, error) // should true, nil to terminate resolving
	}
)

func (r Resolver) Resolve(expr Expr, scope map[string]Def) (Expr, error) {
	return r.resolve(expr, scope, map[string]Def{}, nil)
}

// resolve turn one expression into another where all references points to native types.
// It's a recursive process where each step starts with validation. Invalid expression always leads to error.
// For inst expr it checks compatibility between args and params and returns error if some constraint isn't satisfied.
// Then it updates scope by adding params of ref type with resolved args as values to allow substitution later.
// Then it checks whether base type of current ref type is native type to terminate with nil err and resolved expr.
// For non-native types process starts from the beginning with updated scope. New scope will contain values for params.
// For lit exprs logic is the following: for enum do nothing (it's valid and not composite, there's nothing to resolve),
// for array resolve it's type, for record and union apply recursion for it's every field/element.
func (r Resolver) resolve( //nolint:funlen
	expr Expr,
	scope map[string]Def,
	frame map[string]Def,
	trace *Trace,
) (Expr, error) {
	if err := r.validator.Validate(expr); err != nil { // todo remove embedding
		return Expr{}, fmt.Errorf("%w: %v", ErrInvalidExpr, err)
	}

	switch expr.Lit.Type() {
	case EnumLitType:
		return expr, nil // nothing to resolve in enum
	case ArrLitType:
		resolvedArrType, err := r.resolve(expr.Lit.Arr.Expr, scope, frame, trace)
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
			resolvedEl, err := r.resolve(unionEl, scope, frame, trace)
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
			resolvedFieldExpr, err := r.resolve(fieldExpr, scope, frame, trace)
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

	if def.RecursionAllowed && !def.Body.Empty() {
		return Expr{}, fmt.Errorf("%w: %v", ErrNotBaseTypeSupportsRecursion, def)
	}

	if len(def.Params) != len(expr.Inst.Args) { // args must not be > than params to avoid bad case with constraint
		return Expr{}, fmt.Errorf(
			"%w, want %d, got %d", ErrInstArgsLen, len(def.Params), len(expr.Inst.Args),
		)
	}

	newTrace := Trace{
		prev: trace,
		ref:  expr.Inst.Ref,
	}

	shouldReturn, err := r.terminator.ShouldTerminate(newTrace, scope)
	if err != nil {
		return Expr{}, fmt.Errorf("%w: %v", ErrTerminator, err)
	} else if shouldReturn {
		return expr, nil // IDEA: replace recursive ref with something like `any` (like chat GPT suggested)
	}

	newFrame := make(map[string]Def, len(def.Params))
	resolvedArgs := make([]Expr, 0, len(expr.Inst.Args))
	for i, param := range def.Params { // resolve args and constrs and check their compatibility
		resolvedArg, err := r.resolve(expr.Inst.Args[i], scope, frame, &newTrace)
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrUnresolvedArg, err)
		}
		newFrame[param.Name] = Def{Body: resolvedArg} // no params for generics
		resolvedArgs = append(resolvedArgs, resolvedArg)
		if param.Constr.Empty() {
			continue
		}
		resolvedConstr, err := r.resolve(param.Constr, scope, newFrame, &newTrace) //nolint:lll // we pass newFrame because constr can refer type param
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrConstr, err)
		}

		if err := r.comparator.Check(resolvedArg, newTrace, resolvedConstr, newTrace, scope); err != nil {
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

	return r.resolve(def.Body, scope, newFrame, &newTrace) // TODO replace "flat" args with resolved args?
}

// getDef checks for def in args, then in scope and returns err if expr refers no nothing.
func (Resolver) getDef(ref string, args, scope map[string]Def) (Def, error) {
	def, exist := args[ref]
	if exist {
		return def, nil
	}

	def, exist = scope[ref]
	if !exist {
		return Def{}, fmt.Errorf("%w: %v", ErrUndefinedRef, ref)
	}

	return def, nil
}

func NewDefaultResolver() Resolver {
	return Resolver{
		validator:  Validator{},
		comparator: NewDefaultSubtypeChecker(),
		terminator: RecursionTerminator{},
	}
}

func MustNewResolver(v exprValidator, c compatChecker, t recursionTerminator) Resolver {
	tools.NilPanic(v, c)
	return Resolver{v, c, t}
}
