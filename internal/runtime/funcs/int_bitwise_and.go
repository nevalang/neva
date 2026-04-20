package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intBitwiseAnd struct{}

func (intBitwiseAnd) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	return createBinaryFuncConcurrent(io, func(left runtime.Msg, right runtime.Msg) runtime.Msg {
		return runtime.NewIntMsg(left.Int() & right.Int())
	})
}
