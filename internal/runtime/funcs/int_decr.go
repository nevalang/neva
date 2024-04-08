package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intDecr struct{}

func (i intDecr) Create(io runtime.FuncIO, _ runtime.Msg) (func(context.Context), error) {
	nIn, err := io.In.Port("n")
	if err != nil {
		return nil, err
	}

	nOut, err := io.Out.Port("n")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var n runtime.Msg

		for {
			select {
			case <-ctx.Done():
				return
			case n = <-nIn:
			}

			select {
			case <-ctx.Done():
				return
			case nOut <- runtime.NewIntMsg(n.Int() - 1):
			}
		}
	}, nil
}
