package types

import (
	"errors"
	"fmt"

	"golang.org/x/exp/maps"
)

func (expr Expr) Resolve(scope map[string]Def) (Expr, error) { //nolint:funlen
	if err := expr.Validate(); err != nil {
		return Expr{}, fmt.Errorf("invalid expr: %w", err)
	}

	switch { // resolve literal
	case expr.Lit.EnumLit != nil:
		return expr, nil
	case expr.Lit.UnionLit != nil:
		resolvedUnion := make([]Expr, 0, len(expr.Lit.UnionLit))
		for _, unionEl := range expr.Lit.UnionLit {
			resolvedEl, err := unionEl.Resolve(scope)
			if err != nil {
				return Expr{}, err
			}
			resolvedUnion = append(resolvedUnion, resolvedEl)
		}
		return Expr{
			Lit: LiteralExpr{UnionLit: resolvedUnion},
		}, nil
	case expr.Lit.RecLit != nil:
		resolvedStruct := make(map[string]Expr, len(expr.Lit.RecLit))
		for field, fieldExpr := range expr.Lit.RecLit {
			resolvedFieldExpr, err := fieldExpr.Resolve(scope)
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
		resolvedArg, err := expr.Inst.Args[i].Resolve(scope)
		if err != nil {
			return Expr{}, errors.New("")
		}

		// ОСТОРОЖНО - констрейнт тоже надо ресолвить
		if err := resolvedArg.IsSubTypeOf(param.Constraint); err != nil { // compatibility check
			return Expr{}, errors.New("!resolvedArg.IsSubType")
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
				Inst: InstantiationExpr{
					Ref:  refType.Body.Inst.Ref,
					Args: resolvedArgs,
				},
			}, nil
		}
	}

	return refType.Body.Resolve(newScope) // if it's not a native type and not literal, then do recursive
}
