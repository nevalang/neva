package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type ifRouter struct{}

func (ifRouter) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
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
			msg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			if msg.Bool() {
				if !thenOut.Send(ctx, emptyStruct()) {
					return
				}
			} else {
				if !elseOut.Send(ctx, emptyStruct()) {
					return
				}
			}
		}
	}, nil
}
