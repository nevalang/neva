package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type floatNeg struct{}

func (floatNeg) Create(io runtime.IO, _ messages.Msg) (func(context.Context), error) {
	return createUnaryFunc(io, func(input messages.Msg) messages.Msg {
		return messages.NewFloatMsg(-input.Float())
	})
}
