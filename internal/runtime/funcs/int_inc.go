package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intInc struct{}

func (intInc) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	return createUnaryFunc(io, func(input runtime.Msg) runtime.Msg {
		return runtime.NewIntMsg(input.Int() + 1)
	})
}
