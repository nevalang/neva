package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type not struct{}

func (p not) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}
	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(
				ctx,
				runtime.NewBoolMsg(
					!dataMsg.Bool(),
				),
			) {
				return
			}
		}
	}, nil
}
