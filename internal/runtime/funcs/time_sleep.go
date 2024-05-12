package funcs

import (
	"context"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

type timeSleep struct{}

func (timeSleep) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	durIn, err := io.In.Port("dur")
	if err != nil {
		return nil, err
	}

	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Port("data")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var durMsg runtime.Msg
			select {
			case <-ctx.Done():
				return
			case durMsg = <-durIn:
			}

			var dataMsg runtime.Msg
			select {
			case <-ctx.Done():
				return
			case dataMsg = <-dataIn:
			}

			time.Sleep(time.Duration(durMsg.Int()))

			select {
			case <-ctx.Done():
				return
			case dataOut <- dataMsg:
			}
		}
	}, nil
}
