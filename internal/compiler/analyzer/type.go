package analyzer

import (
	"fmt"

	ts "github.com/nevalang/neva/pkg/typesystem"
)

func (a Analyzer) analyzeTypeDef(def ts.Def, scope Scope) (ts.Def, error) {
	resolvedDef, err := a.resolver.ResolveDef(def, scope)
	if err != nil {
		return ts.Def{}, fmt.Errorf("resolve def: %w", err)
	}
	return resolvedDef, nil
}

func (a Analyzer) analyzeTypeExpr(expr ts.Expr, scope Scope) (ts.Expr, error) {
	resolvedExpr, err := a.resolver.ResolveExpr(expr, scope)
	if err != nil {
		return ts.Expr{}, fmt.Errorf("resolve expr: %w", err)
	}
	return resolvedExpr, nil
}

func (a Analyzer) analyzeTypeParams(params []ts.Param, scope Scope) ([]ts.Param, error) {
	resolvedParams, _, err := a.resolver.ResolveParams(params, scope)
	if err != nil {
		return nil, fmt.Errorf("resolve params: %w", err)
	}
	return resolvedParams, nil
}
