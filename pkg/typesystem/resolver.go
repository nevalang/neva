package typesystem

import (
	"errors"
	"fmt"
)

// Resolver transforms expression it into a form where all references it contains points to resolved expressions.
type Resolver struct {
	validator  exprValidator       // Check if expression invalid before resolving it
	comparator subtypeChecker      // Compare arguments with constraints
	terminator recursionTerminator // Don't stuck in a loop
}

//go:generate mockgen -source $GOFILE -destination mocks_test.go -package ${GOPACKAGE}_test
type (
	exprValidator interface {
		Validate(Expr) error
		ValidateDef(def Def) error
	}
	subtypeChecker interface {
		Check(Expr, Expr, TerminatorParams) error
	}
	recursionTerminator interface {
		ShouldTerminate(Trace, Scope) (bool, error)
	}
	Scope interface {
		GetType(ref string) (Def, error)
		Rebase(ref string) (Scope, error)
	}
)

func (r Resolver) ResolveExpr(expr Expr, scope Scope) (Expr, error) {
	return r.resolveExpr(expr, scope, map[string]Def{}, nil)
}

func (r Resolver) ResolveDef(def Def, scope Scope) (Def, error) {
	resolvedParams, frame, err := r.ResolveParams(def.Params, scope)
	if err != nil {
		return Def{}, fmt.Errorf("resolve params: %w", err)
	}
	if def.BodyExpr == nil {
		return Def{
			Params:                           resolvedParams,
			CanBeUsedForRecursiveDefinitions: def.CanBeUsedForRecursiveDefinitions,
		}, nil
	}
	resolvedBody, err := r.resolveExpr(*def.BodyExpr, scope, frame, nil)
	if err != nil {
		return Def{}, fmt.Errorf("resolve expr: %w", err)
	}
	return Def{
		Params:                           resolvedParams,
		BodyExpr:                         &resolvedBody,
		CanBeUsedForRecursiveDefinitions: def.CanBeUsedForRecursiveDefinitions,
	}, nil
}

func (r Resolver) ResolveParams(params []Param, scope Scope) ([]Param, map[string]Def, error) {
	result := make([]Param, 0, len(params))
	frame := make(map[string]Def, len(params))
	for _, param := range params {
		if param.Constr == nil {
			result = append(result, Param{Name: param.Name})
			continue
		}
		resolved, err := r.resolveExpr(*param.Constr, scope, frame, nil)
		if err != nil {
			return nil, frame, fmt.Errorf("resolve expr: %w", err)
		}
		frame[param.Name] = Def{BodyExpr: &resolved}
		result = append(result, Param{
			Name:   param.Name,
			Constr: &resolved,
		})
	}
	return result, frame, nil
}

// ResolveArgs resolves arguments and parameters and checks that they are compatible.
// It's copy-paste from resolveExpr method because it's hard to reuse that code without creating useless expr and scope.
func (r Resolver) ResolveArgs(args []Expr, params []Param, scope Scope) ([]Expr, []Param, error) {
	newFrame := make(map[string]Def, len(params))
	resolvedArgs := make([]Expr, 0, len(args))
	resolvedParams := make([]Param, 0, len(params))
	for i, param := range params { // resolve args and constrs and check their compatibility
		resolvedArg, err := r.resolveExpr(args[i], scope, nil, nil)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %v", ErrUnresolvedArg, err)
		}

		newFrame[param.Name] = Def{BodyExpr: &resolvedArg} // no params for generics
		resolvedArgs = append(resolvedArgs, resolvedArg)

		if param.Constr != nil {
			resolvedParams = append(resolvedParams, Param{
				Name: param.Name,
			})
			continue
		}

		resolvedConstr, err := r.resolveExpr(*param.Constr, scope, newFrame, nil)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %v", ErrConstr, err)
		}
		resolvedParams = append(resolvedParams, Param{
			Name:   param.Name,
			Constr: &resolvedConstr,
		})

		if err := r.comparator.Check(resolvedArg, resolvedConstr, TerminatorParams{Scope: scope}); err != nil {
			return nil, nil, fmt.Errorf(" %w: %v", ErrIncompatArg, err)
		}
	}
	return resolvedArgs, resolvedParams, nil
}

var (
	ErrInvalidExpr        = errors.New("expression must be valid in order to be resolved")
	ErrScope              = errors.New("can't get type def from scope by ref")
	ErrScopeUpdate        = errors.New("scope update")
	ErrInstArgsLen        = errors.New("inst must have same number of arguments as def has parameters")
	ErrIncompatArg        = errors.New("argument is not subtype of the parameter's contraint")
	ErrUnresolvedArg      = errors.New("can't resolve argument")
	ErrConstr             = errors.New("can't resolve constraint")
	ErrArrType            = errors.New("could not resolve array type")
	ErrUnionUnresolvedEl  = errors.New("can't resolve union element")
	ErrRecFieldUnresolved = errors.New("can't resolve record field")
	ErrInvalidDef         = errors.New("invalid definition")
	ErrTerminator         = errors.New("recursion terminator")
)

// resolveExpr turn one expression into another where all references points to native types.
// It's a recursive process where each step starts with validation. Invalid expression always leads to error.
// For inst expr it checks compatibility between args and params and returns error if some constraint isn't satisfied.
// Then it updates scope by adding params of ref type with resolved args as values to allow substitution later.
// Then it checks whether base type of current ref type is native type to terminate with nil err and resolved expr.
// For non-native types process starts from the beginning with updated scope. New scope will contain values for params.
// For lit exprs logic is the this: for enum do nothing (it's valid and not composite, there's nothing to resolveExpr),
// for array resolveExpr it's type, for record and union apply recursion for it's every field/element.
func (r Resolver) resolveExpr( //nolint:funlen
	expr Expr,
	scope Scope,
	frame map[string]Def,
	trace *Trace,
) (Expr, error) {
	if err := r.validator.Validate(expr); err != nil {
		return Expr{}, fmt.Errorf("%w: %v", ErrInvalidExpr, err)
	}

	switch expr.Lit.Type() {
	case EnumLitType:
		return expr, nil
	case ArrLitType:
		resolvedArrType, err := r.resolveExpr(expr.Lit.Arr.Expr, scope, frame, trace)
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrArrType, err)
		}
		return Expr{
			Lit: &LitExpr{
				Arr: &ArrLit{resolvedArrType, expr.Lit.Arr.Size},
			},
		}, nil
	case UnionLitType:
		resolvedUnion := make([]Expr, 0, len(expr.Lit.Union))
		for _, unionEl := range expr.Lit.Union {
			resolvedEl, err := r.resolveExpr(unionEl, scope, frame, trace)
			if err != nil {
				return Expr{}, fmt.Errorf("%w: %v", ErrUnionUnresolvedEl, err)
			}
			resolvedUnion = append(resolvedUnion, resolvedEl)
		}
		return Expr{
			Lit: &LitExpr{Union: resolvedUnion},
		}, nil
	case RecLitType:
		resolvedStruct := make(map[string]Expr, len(expr.Lit.Rec))
		for field, fieldExpr := range expr.Lit.Rec {
			resolvedFieldExpr, err := r.resolveExpr(fieldExpr, scope, frame, trace)
			if err != nil {
				return Expr{}, fmt.Errorf("%w: %v", ErrRecFieldUnresolved, err)
			}
			resolvedStruct[field] = resolvedFieldExpr
		}
		return Expr{
			Lit: &LitExpr{Rec: resolvedStruct},
		}, nil
	}

	def, scope, err := r.getDef(expr.Inst.Ref, frame, scope)
	if err != nil {
		return Expr{}, err
	}

	if err := r.validator.ValidateDef(def); err != nil {
		return Expr{}, errors.Join(ErrInvalidDef, err)
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
		return expr, nil
	}

	newFrame := make(map[string]Def, len(def.Params))
	resolvedArgs := make([]Expr, 0, len(expr.Inst.Args))
	for i, param := range def.Params { // resolve args and constrs and check their compatibility
		resolvedArg, err := r.resolveExpr(expr.Inst.Args[i], scope, frame, &newTrace)
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrUnresolvedArg, err)
		}

		newFrame[param.Name] = Def{BodyExpr: &resolvedArg} // no params for generics
		resolvedArgs = append(resolvedArgs, resolvedArg)

		if param.Constr != nil {
			continue
		}

		// we pass newFrame because constr can refer type param
		resolvedConstr, err := r.resolveExpr(*param.Constr, scope, newFrame, &newTrace)
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrConstr, err)
		}

		params := TerminatorParams{
			Scope:          scope,
			SubtypeTrace:   newTrace,
			SupertypeTrace: newTrace,
		}

		if err := r.comparator.Check(resolvedArg, resolvedConstr, params); err != nil {
			return Expr{}, fmt.Errorf(" %w: %v", ErrIncompatArg, err)
		}
	}

	if def.BodyExpr == nil {
		return Expr{
			Inst: &InstExpr{
				Ref:  expr.Inst.Ref,
				Args: resolvedArgs,
			},
		}, nil
	}

	return r.resolveExpr(*def.BodyExpr, scope, newFrame, &newTrace)
}

func (Resolver) getDef(ref string, frame map[string]Def, scope Scope) (Def, Scope, error) {
	def, exist := frame[ref]
	if exist {
		return def, scope, nil
	}

	def, err := scope.GetType(ref)
	if err != nil {
		return Def{}, nil, fmt.Errorf("%w: %v", ErrScope, err)
	}

	scope, err = scope.Rebase(ref)
	if err != nil {
		return Def{}, nil, fmt.Errorf("%w: %v", ErrScopeUpdate, err)
	}

	return def, scope, nil
}

func MustNewResolver(validator exprValidator, checker subtypeChecker, terminator recursionTerminator) Resolver {
	if validator == nil || checker == nil || terminator == nil {
		panic("all arguments must be not nil")
	}
	return Resolver{validator, checker, terminator}
}
