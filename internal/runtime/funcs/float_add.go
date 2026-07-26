package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type floatAdd struct{}

func (floatAdd) Create(io runtime.IO, _ messages.Msg) (func(context.Context), error) {
	return createBinaryFuncConcurrent(io, func(left messages.Msg, right messages.Msg) messages.Msg {
		return messages.FloatAdd(left, right)
	})
}
