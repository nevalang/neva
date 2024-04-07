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
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-nIn:
				select {
				case <-ctx.Done():
					return
				case nOut <- runtime.NewIntMsg(data.Int() - 1):
				}
			}
		}
	}, nil
}
