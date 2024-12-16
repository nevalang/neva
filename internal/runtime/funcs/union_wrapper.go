package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type UnionWrapper struct{}

func (UnionWrapper) Create(io runtime.IO, cfg runtime.Msg) (func(ctx context.Context), error) {
	tag := cfg.Str()

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

			if !resOut.Send(ctx, runtime.NewUnionMsg(tag, data)) {
				return
			}
		}
	}, nil
}
