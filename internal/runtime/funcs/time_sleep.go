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

	sigOut, err := io.Out.Port("sig")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var nsMsg runtime.Msg

		for {
			select {
			case <-ctx.Done():
				return
			case nsMsg = <-nsIn:
			}

			time.Sleep(time.Duration(nsMsg.Int()))

			select {
			case <-ctx.Done():
				return
			case sigOut <- nsMsg:
			}
		}
	}, nil
}
