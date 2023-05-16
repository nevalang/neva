package analyzer

import (
	"errors"

	"github.com/nevalang/neva/internal/compiler"
)

func (a Analyzer) analyzeInterface(
	interf compiler.Interface,
	scope Scope,
) (
	compiler.Interface,
	map[compiler.EntityRef]struct{},
	error,
) {
	resolvedParams, used, err := a.analyzeTypeParameters(interf.Params, scope)
	if err != nil {
		return compiler.Interface{}, nil, errors.Join(ErrTypeParams, err)
	}

	resolvedIO, usedByIO, err := a.analyzeIO(interf.IO, scope, resolvedParams)
	if err != nil {
		return compiler.Interface{}, nil, errors.Join(ErrIO, err)
	}

	return compiler.Interface{
		Params: resolvedParams,
		IO:     resolvedIO,
	}, a.mergeUsed(used, usedByIO), nil
}
