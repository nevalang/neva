package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intNeg struct{}

func (intNeg) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	return createUnaryFunc(io, func(input runtime.Msg) runtime.Msg {
		return runtime.NewIntMsg(-input.Int())
	})
}
