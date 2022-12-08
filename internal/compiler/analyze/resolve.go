package analyze

import (
	"github.com/emil14/neva/internal/compiler/src"
)

type Resolver struct{}

func (r Resolver) Resolve(expr src.TypeExpr, scope TypeScope) (src.TypeExpr, error) {
	return resolve(expr, TypeScope{}, scope)
}

func resolve(expr src.TypeExpr, lscope, gscope TypeScope) (src.TypeExpr, error) {
	var refType src.Type

	if expr.Ref.Pkg == "" {
		if _, ok := gscope[expr.Ref]; !ok {
			panic("")
		}
		refType = lscope[expr.Ref]
	}

	if len(refType.Params) != len(expr.RefArgs) {
		panic("len")
	}

	if refType.Body == nil {
		return expr, nil
	}

	newExprCtx := TypeScope{}

	for _, name := range refType.Params {
		ref := src.NewLocalTypeRef(name)
		newExprCtx[ref] = lscope[ref]
	}

	return resolve(*refType.Body, newExprCtx, gscope)
}
