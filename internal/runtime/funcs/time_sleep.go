package funcs

import (
	"context"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

type timeSleep struct{}

func (timeSleep) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	durIn, err := io.In.Single("dur")
	if err != nil {
		return nil, err
	}

	sigOut, err := io.Out.Single("sig")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			durMsg, ok := durIn.Receive(ctx)
			if !ok {
				return
			}

			time.Sleep(time.Duration(durMsg.Int()))

			if !sigOut.Send(ctx, nil) {
				return
			}
		}
	}, nil
}
