package analyzer

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/pkg/sourcecode"
	ts "github.com/nevalang/neva/pkg/typesystem"
)

type analyzeTypeDefParams struct {
	allowEmptyBody bool
}

var ErrEmptyTypeDefBody = fmt.Errorf("Type definition must have non-empty body")

func (a Analyzer) analyzeTypeDef(def ts.Def, scope src.Scope, params analyzeTypeDefParams) (ts.Def, *compiler.Error) {
	if !params.allowEmptyBody && def.BodyExpr == nil {
		meta := def.Meta.(src.Meta) //nolint:forcetypeassert
		return ts.Def{}, &compiler.Error{
			Err:      ErrEmptyTypeDefBody,
			Location: &scope.Location,
			Meta:     &meta,
		}
	}

	// Note that we only resolve params. Body is resolved each time there's an expression that refers to it.
	// We can't resolve body without args. And don't worry about unused bodies. Unused entities are error themselves.
	resolvedParams, err := a.resolver.ResolveParams(def.Params, scope)
	if err != nil {
		meta := def.Meta.(src.Meta) //nolint:forcetypeassert
		return ts.Def{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &meta,
		}
	}

	return ts.Def{
		Params:                           resolvedParams,
		BodyExpr:                         def.BodyExpr,
		CanBeUsedForRecursiveDefinitions: def.CanBeUsedForRecursiveDefinitions,
	}, nil
}

func (a Analyzer) analyzeTypeExpr(expr ts.Expr, scope src.Scope) (ts.Expr, *compiler.Error) {
	resolvedExpr, err := a.resolver.ResolveExpr(expr, scope)
	if err != nil {
		meta := expr.Meta.(src.Meta) //nolint:forcetypeassert
		return ts.Expr{}, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
			Meta:     &meta,
		}
	}
	return resolvedExpr, nil
}

func (a Analyzer) analyzeTypeParams(params []ts.Param, scope src.Scope) ([]ts.Param, *compiler.Error) {
	resolvedParams, err := a.resolver.ResolveParams(params, scope)
	if err != nil {
		return nil, &compiler.Error{
			Err:      err,
			Location: &scope.Location,
		}
	}
	return resolvedParams, nil
}
