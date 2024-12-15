package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type new struct{}

func (c new) Create(io runtime.IO, cfg runtime.Msg) (func(ctx context.Context), error) {
	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if !resOut.Send(ctx, cfg) {
				return
			}
		}
	}, nil
}
