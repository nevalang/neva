package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type if_ struct{}

func (p if_) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	thenOut, err := io.Out.Single("then")
	if err != nil {
		return nil, err
	}

	elseOut, err := io.Out.Single("else")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			var out runtime.SingleOutport
			if dataMsg.Bool() {
				out = thenOut
			} else {
				out = elseOut
			}

			if !out.Send(ctx, dataMsg) {
				return
			}
		}
	}, nil
}
