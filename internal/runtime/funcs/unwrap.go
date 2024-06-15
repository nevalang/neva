package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type unwrap struct{}

func (unwrap) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	someOut, err := io.Out.SingleOutport("some")
	if err != nil {
		return nil, err
	}

	noneOut, err := io.Out.SingleOutport("none")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			if dataMsg == nil {
				if !noneOut.Send(ctx, nil) {
					return
				}
				continue
			}

			if !someOut.Send(ctx, dataMsg) {
				return
			}
		}
	}, nil
}
