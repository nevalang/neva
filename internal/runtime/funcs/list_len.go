package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type listlen struct{}

func (p listlen) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.SingleInport("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.SingleOutport("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var data runtime.Msg

		for {
			select {
			case <-ctx.Done():
				return
			case data = <-dataIn:
			}

			l := len(data.List())

			select {
			case <-ctx.Done():
				return
			case resOut <- runtime.NewIntMsg(int64(l)):
			}
		}
	}, nil
}
