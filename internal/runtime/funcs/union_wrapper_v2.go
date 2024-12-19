package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

// UnionWrapV2 wraps with extra signal
type UnionWrapV2 struct{}

func (UnionWrapV2) Create(io runtime.IO, cfg runtime.Msg) (func(ctx context.Context), error) {
	tag := cfg.Str()

	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	signalIn, err := io.In.Single("sig")
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

			if _, ok := signalIn.Receive(ctx); !ok {
				return
			}

			if !resOut.Send(ctx, runtime.NewUnionMsg(tag, data)) {
				return
			}
		}
	}, nil
}
