package funcs

import (
	"context"
	"slices"

	"github.com/nevalang/neva/internal/runtime"
)

type listPush struct{}

func (p listPush) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}
	lstIn, err := io.In.Single("lst")
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

			lstMsg, ok := lstIn.Receive(ctx)
			if !ok {
				return
			}

			lstCopy := slices.Clone(lstMsg.List())

			if !resOut.Send(
				ctx,
				runtime.NewListMsg(
					append(lstCopy, dataMsg),
				),
			) {
				return
			}
		}
	}, nil
}
