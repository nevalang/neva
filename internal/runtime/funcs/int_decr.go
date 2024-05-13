package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intDecr struct{}

func (i intDecr) Create(io runtime.FuncIO, _ runtime.Msg) (func(context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var dataMsg runtime.Msg

		for {
			select {
			case <-ctx.Done():
				return
			case dataMsg = <-dataIn:
			}

			select {
			case <-ctx.Done():
				return
			case resOut <- runtime.NewIntMsg(dataMsg.Int() - 1):
			}
		}
	}, nil
}
