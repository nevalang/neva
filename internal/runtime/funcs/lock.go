package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type lock struct{}

func (l lock) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := io.In.Single("sig")
	if err != nil {
		return nil, err
	}

	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.SingleOutport("data")
	if err != nil {
		return nil, err
	}

	return l.Handle(sigIn, dataIn, dataOut), nil
}

func (lock) Handle(
	sigIn runtime.SingleInport,
	dataIn runtime.SingleInport,
	dataOut runtime.SingleOutport,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			if !dataOut.Send(ctx, data) {
				return
			}
		}
	}
}
