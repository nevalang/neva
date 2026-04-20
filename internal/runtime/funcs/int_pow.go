package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intPow struct{}

func (intPow) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	return createBinaryFuncConcurrent(io, func(left runtime.Msg, right runtime.Msg) runtime.Msg {
		base := left.Int()
		exponent := right.Int()
		result := int64(1)

		for range exponent {
			result *= base
		}

		return runtime.NewIntMsg(result)
	})
}
