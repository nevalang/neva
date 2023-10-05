package analyzer

import (
	"fmt"

	ts "github.com/nevalang/neva/pkg/typesystem"
)

type analyzeTypeDefParams struct {
	allowEmptyBody bool
}

var ErrEmptyTypeDefBody = fmt.Errorf("type def body is empty")

func (a Analyzer) analyzeTypeDef(def ts.Def, scope Scope, params analyzeTypeDefParams) (ts.Def, error) {
	if !params.allowEmptyBody && def.BodyExpr == nil {
		return ts.Def{}, ErrEmptyTypeDefBody
	}

	// Note that we only resolve params. Body is resolved each time there's an expression that refers to it.
	// We can't resolve body without args. And don't worry about unused bodies. Unused entities are error themselves.
	resolvedParams, err := a.resolver.ResolveParams(def.Params, scope)
	if err != nil {
		return ts.Def{}, fmt.Errorf("resolve def: %w", err)
	}

	return ts.Def{
		Params:                           resolvedParams,
		BodyExpr:                         def.BodyExpr,
		CanBeUsedForRecursiveDefinitions: def.CanBeUsedForRecursiveDefinitions,
	}, nil
}

func (a Analyzer) analyzeTypeExpr(expr ts.Expr, scope Scope) (ts.Expr, error) {
	resolvedExpr, err := a.resolver.ResolveExpr(expr, scope)
	if err != nil {
		return ts.Expr{}, fmt.Errorf("resolve expr: %w", err)
	}
	return resolvedExpr, nil
}

func (a Analyzer) analyzeTypeParams(params []ts.Param, scope Scope) ([]ts.Param, error) {
	resolvedParams, err := a.resolver.ResolveParams(params, scope)
	if err != nil {
		return nil, fmt.Errorf("resolve params: %w", err)
	}
	return resolvedParams, nil
}
