package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type new struct{}

func (c new) Create(io runtime.FuncIO, msg runtime.Msg) (func(ctx context.Context), error) {
	outport, err := io.Out.SingleOutport("msg")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if !outport.Send(ctx, msg) {
				return
			}
		}
	}, nil
}
