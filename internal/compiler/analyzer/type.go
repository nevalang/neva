package analyzer

import (
	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	ts "github.com/nevalang/neva/internal/compiler/sourcecode/typesystem"
)

type analyzeTypeDefParams struct {
	allowEmptyBody bool
}

func (a Analyzer) analyzeType(def ts.Def, scope src.Scope, params analyzeTypeDefParams) (ts.Def, *compiler.Error) {
	if !params.allowEmptyBody && def.BodyExpr == nil {
		meta := def.Meta
		return ts.Def{}, &compiler.Error{
			Message: "Type definition must have non-empty body",
			Meta:    &meta,
		}
	}

	// Note that we only resolve params. Body is resolved each time there's an expression that refers to it.
	// We can't resolve body without args. And don't worry about unused bodies. Unused entities are error themselves.
	resolvedParams, _, err := a.resolver.ResolveParams(def.Params, scope)
	if err != nil {
		meta := def.Meta
		return ts.Def{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &meta,
		}
	}

	return ts.Def{
		Params:   resolvedParams,
		BodyExpr: def.BodyExpr,
	}, nil
}

func (a Analyzer) analyzeTypeExpr(expr ts.Expr, scope src.Scope) (ts.Expr, *compiler.Error) {
	resolvedExpr, err := a.resolver.ResolveExpr(expr, scope)
	if err != nil {
		meta := expr.Meta //nolint:forcetypeassert
		return ts.Expr{}, &compiler.Error{
			Message: err.Error(),
			Meta:    &meta,
		}
	}
	return resolvedExpr, nil
}

func (a Analyzer) analyzeTypeParams(
	params []ts.Param,
	scope src.Scope,
) (
	[]ts.Param,
	*compiler.Error,
) {
	resolvedParams, _, err := a.resolver.ResolveParams(params, scope)
	if err != nil {
		return nil, &compiler.Error{
			Message: err.Error(),
		}
	}
	return resolvedParams, nil
}
