package analyzer

import (
	"fmt"

	ts "github.com/nevalang/neva/pkg/typesystem"
)

func (a Analyzer) analyzeTypeDef(def ts.Def) (ts.Def, error) {
	resolvedDef, err := a.resolver.ResolveDef(def, nil)
	if err != nil {
		return ts.Def{}, fmt.Errorf("resolve def: %w", err)
	}
	return resolvedDef, nil
}

func (a Analyzer) analyzeTypeExpr(expr ts.Expr) (ts.Expr, error) {
	resolvedExpr, err := a.resolver.ResolveExpr(expr, nil)
	if err != nil {
		return ts.Expr{}, fmt.Errorf("resolve expr: %w", err)
	}
	return resolvedExpr, nil
}

func (a Analyzer) analyzeTypeParams(params []ts.Param) ([]ts.Param, error) {
	resolvedParams, _, err := a.resolver.ResolveParams(params, nil)
	if err != nil {
		return nil, fmt.Errorf("resolve params: %w", err)
	}
	return resolvedParams, nil
}
