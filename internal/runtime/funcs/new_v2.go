package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type newV2 struct{}

func (c newV2) Create(io runtime.IO, cfg runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := io.In.Single("sig")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}
			if !resOut.Send(ctx, cfg) {
				return
			}
		}
	}, nil
}
