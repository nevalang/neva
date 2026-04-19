package funcs

import (
	"context"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

type timeDelay struct{}

//nolint:varnamelen
func (timeDelay) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	durIn, err := io.In.Single("dur")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			//nolint:varnamelen
			durMsg, ok := durIn.Receive(ctx)
			if !ok {
				return
			}

			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			time.Sleep(time.Duration(durMsg.Int()))

			if !resOut.Send(ctx, dataMsg) {
				return
			}
		}
	}, nil
}
