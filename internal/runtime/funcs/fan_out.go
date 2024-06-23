package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type fanOut struct{}

func (d fanOut) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, _ := io.In.Single("data")
	dataOut, _ := io.Out.Array("data")

	return func(ctx context.Context) {
		for {
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			if !dataOut.SendAll(ctx, dataMsg) {
				return
			}
		}
	}, nil
}
