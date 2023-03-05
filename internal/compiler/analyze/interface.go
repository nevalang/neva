package analyze

import (
	"errors"

	"github.com/emil14/neva/internal/compiler/src"
	ts "github.com/emil14/neva/pkg/types"
)

func (a Analyzer) analyzeInterface(
	interf src.Interface,
	scope Scope,
	args map[string]ts.Expr,
) (
	src.Interface,
	map[src.EntityRef]struct{},
	error,
) {
	resolvedParams, used, err := a.analyzeTypeParameters(interf.TypeParams, scope, args)
	if err != nil {
		return src.Interface{}, nil, errors.Join(ErrTypeParams, err)
	}

	resolvedIO, usedByIO, err := a.analyzeIO(interf.IO, scope, resolvedParams)
	if err != nil {
		return src.Interface{}, nil, errors.Join(ErrIO, err)
	}

	return src.Interface{
		TypeParams: []ts.Param{},
		IO:         resolvedIO,
	}, a.mergeUsed(used, usedByIO), nil
}
