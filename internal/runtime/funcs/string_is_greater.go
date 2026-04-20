package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type strIsGreater struct{}

func (strIsGreater) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	return createBinaryFuncSequential(io, func(left runtime.Msg, right runtime.Msg) runtime.Msg {
		return runtime.NewBoolMsg(left.Str() > right.Str())
	})
}
