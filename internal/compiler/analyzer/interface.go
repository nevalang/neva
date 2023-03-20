package analyzer

import (
	"errors"

	"github.com/emil14/neva/internal/compiler/src"
)

func (a Analyzer) analyzeInterface(
	interf src.Interface,
	scope Scope,
) (
	src.Interface,
	map[src.EntityRef]struct{},
	error,
) {
	resolvedParams, used, err := a.analyzeTypeParameters(interf.Params, scope)
	if err != nil {
		return src.Interface{}, nil, errors.Join(ErrTypeParams, err)
	}

	resolvedIO, usedByIO, err := a.analyzeIO(interf.IO, scope, resolvedParams)
	if err != nil {
		return src.Interface{}, nil, errors.Join(ErrIO, err)
	}

	return src.Interface{
		Params: resolvedParams,
		IO:     resolvedIO,
	}, a.mergeUsed(used, usedByIO), nil
}
