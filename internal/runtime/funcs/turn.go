package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

// turn serializes data messages after the first message.
type turn struct{}

func (turn) Create(runtimeIO runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	doneIn, err := singleInport(runtimeIO, "done")
	if err != nil {
		return nil, err
	}

	dataIn, err := singleInport(runtimeIO, "data")
	if err != nil {
		return nil, err
	}

	resOut, err := singleOutport(runtimeIO, "res")
	if err != nil {
		return nil, err
	}

	return newTurnHandler(doneIn, dataIn, resOut), nil
}

// newTurnHandler releases the first message, then pairs each following message with done.
func newTurnHandler(
	doneIn, dataIn runtime.SingleInport,
	resOut runtime.SingleOutport,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		first := true
		for {
			if !first {
				if _, received := doneIn.Receive(ctx); !received {
					return
				}
			}

			data, received := dataIn.Receive(ctx)
			if !received {
				return
			}

			if !resOut.Send(ctx, data) {
				return
			}
			first = false
		}
	}
}
