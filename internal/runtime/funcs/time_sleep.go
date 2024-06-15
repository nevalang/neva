package funcs

import (
	"context"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

type timeSleep struct{}

func (timeSleep) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	durIn, err := io.In.Single("dur")
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

	return func(ctx context.Context) {
		for {
			durMsg, ok := durIn.Receive(ctx)
			if !ok {
				return
			}

			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			time.Sleep(time.Duration(durMsg.Int()))

			if !dataOut.Send(ctx, dataMsg) {
				return
			}
		}
	}, nil
}
