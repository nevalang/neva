package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type unionWrapper struct{}

func (unionWrapper) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	tagIn, err := io.In.Single("tag")
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

			tagMsg, ok := tagIn.Receive(ctx)
			if !ok {
				return
			}

			tag := tagMsg.Union().Tag()
			if !resOut.Send(ctx, runtime.NewUnionMsg(tag, dataMsg)) {
				return
			}
		}
	}, nil
}
