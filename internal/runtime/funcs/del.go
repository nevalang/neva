package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type del struct{}

func (d del) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if _, ok := dataIn.Receive(ctx); !ok {
				return
			}
		}
	}, nil
}
