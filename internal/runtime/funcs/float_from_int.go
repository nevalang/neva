package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type floatFromInt struct{}

func (floatFromInt) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(ctx, runtime.NewFloatMsg(float64(data.Int()))) {
				return
			}
		}
	}, nil
}
