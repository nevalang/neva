package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type void struct{}

func (v void) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	vin, err := io.In.Port("v")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-vin:
			}
		}
	}, nil
}
