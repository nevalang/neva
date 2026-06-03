package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type floatSub struct{}

func (floatSub) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	return createBinaryFuncConcurrent(io, func(left runtime.Msg, right runtime.Msg) runtime.Msg {
		return runtime.NewFloatMsg(left.Float() - right.Float())
	})
}
