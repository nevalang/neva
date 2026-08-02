package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

// turn serializes data messages after the first message.
type turn struct{}

func (turn) Create(runtimeIO runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
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

// receiveTurnInputs receives the first data message directly and each following data/done pair concurrently.
func receiveTurnInputs(
	ctx context.Context,
	first bool,
	dataIn, doneIn runtime.SingleInport,
) (runtime.OrderedMsg, runtime.OrderedMsg, bool) {
	if first {
		data, received := dataIn.Receive(ctx)
		return data, runtime.OrderedMsg{}, received
	}

	// The previous acknowledgement and the next data message form one operation.
	// Receive both concurrently so either sender can arrive first.
	return receive2(ctx, dataIn, doneIn)
}

// newTurnHandler releases the first message, then pairs each following message with done.
func newTurnHandler(
	doneIn, dataIn runtime.SingleInport,
	resOut runtime.SingleOutport,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		first := true
		for {
			data, done, received := receiveTurnInputs(ctx, first, dataIn, doneIn)
			if !received {
				return
			}

			if first {
				if !resOut.Send(ctx, data) {
					return
				}
			} else if !resOut.Send(ctx, data, data, done) {
				return
			}

			first = false
		}
	}
}
