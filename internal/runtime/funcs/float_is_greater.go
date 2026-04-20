package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type floatIsGreater struct{}

func (floatIsGreater) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	return createBinaryFuncSequential(io, func(left runtime.Msg, right runtime.Msg) runtime.Msg {
		return runtime.NewBoolMsg(left.Float() > right.Float())
	})
}
