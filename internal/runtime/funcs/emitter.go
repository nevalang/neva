package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type emitter struct{}

func (c emitter) Create(io runtime.FuncIO, msg runtime.Msg) (func(ctx context.Context), error) {
	outport, err := io.Out.Port("msg")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case outport <- msg:
			}
		}
	}, nil
}
