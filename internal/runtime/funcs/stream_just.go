package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamJust struct{}

func (streamJust) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
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
			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}
			if !resOut.Send(ctx, streamOpen()) {
				return
			}
			if !resOut.Send(ctx, streamData(data)) {
				return
			}
			if !resOut.Send(ctx, streamClose()) {
				return
			}
		}
	}, nil
}
