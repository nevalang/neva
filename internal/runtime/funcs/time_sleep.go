package funcs

import (
	"context"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

type timeSleep struct{}

func (timeSleep) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	nsIn, err := io.In.Port("ns")
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
			var nsMsg runtime.Msg
			select {
			case <-ctx.Done():
				return
			case nsMsg = <-nsIn:
			}

			var dataMsg runtime.Msg
			select {
			case <-ctx.Done():
				return
			case dataMsg = <-dataIn:
			}

			time.Sleep(time.Duration(nsMsg.Int()))

			select {
			case <-ctx.Done():
				return
			case dataOut <- dataMsg:
			}
		}
	}, nil
}
