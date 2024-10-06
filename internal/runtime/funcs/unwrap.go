package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type unwrap struct{}

func (unwrap) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	someOut, err := io.Out.Single("some")
	if err != nil {
		return nil, err
	}

	noneOut, err := io.Out.Single("none")
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
				if !noneOut.Send(ctx, emptyStruct()) {
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
