package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intDecrementor struct{}

func (i intDecrementor) Create(io runtime.FuncIO, _ runtime.Msg) (func(context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-dataIn:
				select {
				case <-ctx.Done():
					return
				case resOut <- runtime.NewIntMsg(data.Int() - 1):
				}
			}
		}
	}, nil
}
