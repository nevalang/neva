package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type new struct{}

func (c new) Create(io runtime.FuncIO, cfg runtime.Msg) (func(ctx context.Context), error) {
	outport, err := io.Out.Single("msg")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if !outport.Send(ctx, cfg) {
				return
			}
		}
	}, nil
}
