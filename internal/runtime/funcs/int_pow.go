package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type intPow struct{}

func (intPow) Create(io runtime.IO, _ messages.Msg) (func(context.Context), error) {
	return createBinaryFuncConcurrent(io, func(left messages.Msg, right messages.Msg) messages.Msg {
		base := left.Int()
		exponent := right.Int()
		result := int64(1)

		for range exponent {
			result *= base
		}

		return messages.NewIntMsg(result)
	})
}
