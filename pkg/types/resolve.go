package types

import (
	"errors"
	"fmt"

	"golang.org/x/exp/maps"
)

func (expr Expr) Resolve(scope map[string]Def) (Expr, error) { //nolint:funlen
	switch { // resolve literal
	case expr.Literal.EnumLit != nil:
		return expr, nil
	case expr.Literal.UnionLit != nil:
		resolvedUnion := make([]Expr, 0, len(expr.Literal.UnionLit))
		for _, unionEl := range expr.Literal.UnionLit {
			resolvedEl, err := unionEl.Resolve(scope)
			if err != nil {
				return Expr{}, err
			}
			resolvedUnion = append(resolvedUnion, resolvedEl)
		}
		return Expr{
			Literal: LiteralExpr{UnionLit: resolvedUnion},
		}, nil
	case expr.Literal.RecLit != nil:
		resolvedStruct := make(map[string]Expr, len(expr.Literal.RecLit))
		for field, fieldExpr := range expr.Literal.RecLit {
			resolvedFieldExpr, err := fieldExpr.Resolve(scope)
			if err != nil {
				return Expr{}, errors.New("")
			}
			resolvedStruct[field] = resolvedFieldExpr
		}
		return Expr{
			Literal: LiteralExpr{RecLit: resolvedStruct},
		}, nil
	}

	refType, ok := scope[expr.Instantiation.Ref] // check that reference type exists
	if !ok {
		return Expr{}, errors.New("ref type not found in scope")
	}

	if len(refType.Params) > len(expr.Instantiation.Args) { // check that generic args for every param is present
		return Expr{}, fmt.Errorf(
			"expr must have at least %d arguments, got %d",
			len(refType.Params), len(expr.Instantiation.Args),
		)
	}

	newScope := make(map[string]Def, len(scope)+len(refType.Params)) // new scope contains resolved params (shadow)
	maps.Copy(newScope, scope)
	resolvedArgs := make([]Expr, 0, len(refType.Params))

	for i, param := range refType.Params {
		resolvedArg, err := expr.Instantiation.Args[i].Resolve(scope)
		if err != nil {
			return Expr{}, errors.New("")
		}

		// ОСТОРОЖНО - констрейнт тоже надо ресолвить
		if err := resolvedArg.IsSubType(param.Constraint); err != nil { // compatibility check
			return Expr{}, errors.New("!resolvedArg.IsSubType")
		}

		resolvedArgs = append(resolvedArgs, resolvedArg)
		newScope[param.Name] = Def{
			Params: nil, // we don't refer generics with another generics inside
			Body:   resolvedArg,
		}
	}

	if refType.Body.Literal.Empty() { // reference type's body is an instantiation
		baseType, ok := scope[refType.Body.Instantiation.Ref]
		if !ok {
			return Expr{}, errors.New("")
		}
		if expr.Instantiation.Ref == baseType.Body.Instantiation.Ref { // direct self reference = native instantiation
			return Expr{
				Instantiation: InstantiationExpr{
					Ref:  refType.Body.Instantiation.Ref,
					Args: resolvedArgs,
				},
			}, nil
		}
	}

	return refType.Body.Resolve(newScope) // if it's not a native type and not literal, then do recursive
}
