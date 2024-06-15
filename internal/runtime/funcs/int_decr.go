package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intDecr struct{}

func (i intDecr) Create(io runtime.FuncIO, _ runtime.Msg) (func(context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.SingleOutport("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(ctx, runtime.NewIntMsg(dataMsg.Int()-1)) {
				return
			}
		}
	}, nil
}
