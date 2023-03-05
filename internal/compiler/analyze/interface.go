package analyze

import (
	"errors"

	"github.com/emil14/neva/internal/compiler/src"
	"github.com/emil14/neva/pkg/types"
)

func (a Analyzer) analyzeInterface(interf src.Interface, scope Scope) (src.Interface, map[src.EntityRef]struct{}, error) {
	resolvedParams, used, err := a.analyzeTypeParameters(interf.TypeParams, scope)
	if err != nil {
		return src.Interface{}, nil, errors.Join(ErrTypeParams, err)
	}

	resolvedIO, usedByIO, err := a.analyzeIO(interf.IO, scope, resolvedParams)
	if err != nil {
		return src.Interface{}, nil, errors.Join(ErrIO, err)
	}

	return src.Interface{
		TypeParams: []types.Param{},
		IO:         resolvedIO,
	}, a.mergeUsed(used, usedByIO), nil
}
