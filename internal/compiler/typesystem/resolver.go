package typesystem

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/ast/core"
)

var (
	ErrInvalidExpr        = errors.New("expression must be valid in order to be resolved")
	ErrScope              = errors.New("can't get type def from scope by ref")
	ErrScopeUpdate        = errors.New("scope update")
	ErrInstArgsCount      = errors.New("wrong number of type arguments")
	ErrIncompatArg        = errors.New("argument is not subtype of the parameter's contraint")
	ErrUnresolvedArg      = errors.New("can't resolve argument")
	ErrConstr             = errors.New("can't resolve constraint")
	ErrUnionUnresolvedEl  = errors.New("can't resolve union element")
	ErrRecFieldUnresolved = errors.New("can't resolve struct field")
	ErrInvalidDef         = errors.New("invalid definition")
	ErrTerminator         = errors.New("recursion terminator")
)

// Resolver transforms expression it into a form where all references it contains points to resolved expressions.
type Resolver struct {
	validator  exprValidator       // Check if expression invalid before resolving it
	checker    subtypeChecker      // Compare arguments with constraints
	terminator recursionTerminator // Don't stuck in a loop
}

//go:generate mockgen -source $GOFILE -destination mocks_test.go -package ${GOPACKAGE}_test
type (
	exprValidator interface {
		Validate(expr Expr) error
		ValidateDef(def Def) error
	}
	subtypeChecker interface {
		Check(sub Expr, sup Expr, params TerminatorParams) error
	}
	recursionTerminator interface {
		ShouldTerminate(trace Trace, scope Scope) (bool, error)
	}
	Scope interface {
		GetType(ref core.EntityRef) (Def, Scope, error)
		IsTopType(expr Expr) bool
	}
)

// ResolveExpr resolves given expression using only global scope.
func (r Resolver) ResolveExpr(expr Expr, scope Scope) (Expr, error) {
	return r.resolveExpr(expr, scope, map[string]Def{}, nil)
}

// ResolveExprWithFrame works like ResolveExpr but allows to pass local scope.
func (r Resolver) ResolveExprWithFrame(
	expr Expr,
	frame map[string]Def,
	scope Scope,
) (Expr, error) {
	return r.resolveExpr(expr, scope, frame, nil)
}

// ResolveExprWithFrame works like ResolveExprWithFrame but for list of expressions.
func (r Resolver) ResolveExprsWithFrame(
	exprs []Expr,
	frame map[string]Def,
	scope Scope,
) ([]Expr, error) {
	resolvedExprs := make([]Expr, 0, len(exprs))
	for _, expr := range exprs {
		resolvedExpr, err := r.resolveExpr(expr, scope, frame, nil)
		if err != nil {
			return nil, err
		}
		resolvedExprs = append(resolvedExprs, resolvedExpr)
	}
	return resolvedExprs, nil
}

// ResolveParams resolves every constraint in given parameter list.
func (r Resolver) ResolveParams(
	params []Param,
	scope Scope,
) (
	[]Param, // resolved parameters
	map[string]Def, // resolved frame `paramName:resolvedConstr`
	error,
) {
	result := make([]Param, 0, len(params))
	frame := make(map[string]Def, len(params))
	for _, param := range params {
		resolved, err := r.resolveExpr(param.Constr, scope, frame, nil)
		if err != nil {
			return nil, nil, fmt.Errorf("resolve expr: %w", err)
		}
		frame[param.Name] = Def{BodyExpr: &resolved}
		result = append(result, Param{
			Name:   param.Name,
			Constr: resolved,
		})
	}
	return result, frame, nil
}

// IsSubtypeOf resolves both `sub` and `sup` expressions
// and returns error if `sub` is not subtype of `sup`.
func (r Resolver) IsSubtypeOf(sub, sup Expr, scope Scope) error {
	resolvedSub, err := r.resolveExpr(sub, scope, nil, nil)
	if err != nil {
		return fmt.Errorf("resolve sub expr: %w", err)
	}
	resolvedSup, err := r.resolveExpr(sup, scope, nil, nil)
	if err != nil {
		return fmt.Errorf("resolve sup expr: %w", err)
	}
	return r.checker.Check(
		resolvedSub,
		resolvedSup,
		TerminatorParams{Scope: scope},
	)
}

// CheckArgsCompatibility resolves args
// and params and then checks their compatibility.
func (r Resolver) CheckArgsCompatibility(args []Expr, params []Param, scope Scope) error {
	if len(args) != len(params) {
		return fmt.Errorf(
			"count of arguments mismatch count of parameters, want %d got %d",
			len(params),
			len(args),
		)
	}

	for i := range params {
		arg := args[i]
		param := params[i]

		resolvedSub, err := r.resolveExpr(arg, scope, nil, nil)
		if err != nil {
			return fmt.Errorf("resolve arg expr: %w", err)
		}

		resolvedSup, err := r.resolveExpr(param.Constr, scope, nil, nil)
		if err != nil {
			return fmt.Errorf("resolve param constr expr: %w", err)
		}

		if err := r.checker.Check(
			resolvedSub,
			resolvedSup,
			TerminatorParams{Scope: scope},
		); err != nil {
			return err
		}
	}

	return nil
}

// resolveExpr turn one expression into another where all references points to native types.
// It's a recursive process where each step starts with validation. Invalid expression always leads to error.
// For inst expr it checks compatibility between args and params and returns error if some constraint isn't satisfied.
// Then it updates scope by adding params of ref type with resolved args as values to allow substitution later.
// Then it checks whether base type of current ref type is native type to terminate with nil err and resolved expr.
// For non-native types process starts from the beginning with updated scope. New scope will contain values for params.
// For lit exprs logic is the this:
// for struct and union apply recursion for it's every field/element.
func (r Resolver) resolveExpr(
	expr Expr, // expression to be resolved
	scope Scope, // global scope
	frame map[string]Def, // local scope
	trace *Trace, // how did we get here
) (Expr, error) {
	if err := r.validator.Validate(expr); err != nil {
		return Expr{}, fmt.Errorf("%w: %v", ErrInvalidExpr, err)
	}

	if expr.Lit != nil {
		switch expr.Lit.Type() {
		case UnionLitType:
			resolvedUnion := make(map[string]*Expr, len(expr.Lit.Union))
			for unionElName, unionEl := range expr.Lit.Union {
				if unionEl == nil {
					resolvedUnion[unionElName] = nil
					continue
				}
				resolvedEl, err := r.resolveExpr(*unionEl, scope, frame, trace)
				if err != nil {
					return Expr{}, fmt.Errorf("%w: %v", ErrUnionUnresolvedEl, err)
				}
				resolvedUnion[unionElName] = &resolvedEl
			}
			return Expr{Lit: &LitExpr{Union: resolvedUnion}}, nil
		case StructLitType:
			resolvedStruct := make(map[string]Expr, len(expr.Lit.Struct))
			for field, fieldExpr := range expr.Lit.Struct {
				// we create new trace with virtual ref "struct" (it's safe because it's reserved word)
				// otherwise expressions like `error struct {child maybe<error>}` will be direct recursive for terminator
				newTrace := Trace{
					prev: trace,
					cur:  core.EntityRef{Name: "struct"},
				}
				resolvedFieldExpr, err := r.resolveExpr(
					fieldExpr,
					scope,
					frame,
					&newTrace,
				)
				if err != nil {
					return Expr{}, fmt.Errorf(
						"%w: %v: %v",
						ErrRecFieldUnresolved,
						field,
						err,
					)
				}
				resolvedStruct[field] = resolvedFieldExpr
			}
			return Expr{
				Lit: &LitExpr{Struct: resolvedStruct},
			}, nil
		}
	}

	def, scopeWhereDefFound, err := r.getDef(expr.Inst.Ref, frame, scope)
	if err != nil {
		return Expr{}, err
	}

	if err := r.validator.ValidateDef(def); err != nil {
		return Expr{}, errors.Join(ErrInvalidDef, err)
	}

	if len(def.Params) != len(expr.Inst.Args) { // args must not be > than params to avoid bad case with constraint
		return Expr{}, fmt.Errorf(
			"%w for '%v': want %d, got %d",
			ErrInstArgsCount,
			expr.Inst.Ref,
			len(def.Params),
			len(expr.Inst.Args),
		)
	}

	newTrace := Trace{
		prev: trace,
		cur:  expr.Inst.Ref, // FIXME t1
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

		// we pass newFrame because constr can refer to type parameters
		resolvedConstr, err := r.resolveExpr(
			param.Constr,
			scope,
			newFrame,
			&newTrace,
		)
		if err != nil {
			return Expr{}, fmt.Errorf("%w: %v", ErrConstr, err)
		}

		params := TerminatorParams{
			Scope:          scope,
			SubtypeTrace:   newTrace,
			SupertypeTrace: newTrace,
		}

		if err := r.checker.Check(resolvedArg, resolvedConstr, params); err != nil {
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

	return r.resolveExpr(*def.BodyExpr, scopeWhereDefFound, newFrame, &newTrace)
}

func (Resolver) getDef(
	ref core.EntityRef,
	frame map[string]Def,
	scope Scope,
) (Def, Scope, error) {
	strRef := ref.String()
	def, exist := frame[strRef]
	if exist {
		return def, scope, nil
	}

	def, scope, err := scope.GetType(ref)
	if err != nil {
		return Def{}, nil, fmt.Errorf("%w: %v", ErrScope, err)
	}

	return def, scope, nil
}

func MustNewResolver(validator exprValidator, checker subtypeChecker, terminator recursionTerminator) Resolver {
	if validator == nil || checker == nil || terminator == nil {
		panic("all arguments must be not nil")
	}
	return Resolver{validator, checker, terminator}
}
