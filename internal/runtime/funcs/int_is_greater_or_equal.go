package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intIsGreaterOrEqual struct{}

func (intIsGreaterOrEqual) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	return createBinaryFuncConcurrent(io, func(left runtime.Msg, right runtime.Msg) runtime.Msg {
		return runtime.NewBoolMsg(left.Int() >= right.Int())
	})
}
