package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type fanOut struct{}

func (d fanOut) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Array("data")
	if err != nil {
		return nil, err
	}

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
