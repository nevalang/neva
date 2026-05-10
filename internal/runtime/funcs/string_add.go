package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type stringAdd struct{}

func (stringAdd) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	return createBinaryFuncConcurrent(io, func(left runtime.Msg, right runtime.Msg) runtime.Msg {
		return runtime.NewStringMsg(left.Str() + right.Str())
	})
}
