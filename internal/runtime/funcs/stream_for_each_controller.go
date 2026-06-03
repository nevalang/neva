package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

// streamForEachController serializes stream item handling for ForEach.
// It forwards Open/Close immediately and forwards Data only after done signal.
type streamForEachController struct{}

func (streamForEachController) Create(runtimeIO runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := singleInport(runtimeIO, "data")
	if err != nil {
		return nil, err
	}

	doneIn, err := singleInport(runtimeIO, "done")
	if err != nil {
		return nil, err
	}

	itemOut, err := singleOutport(runtimeIO, "item")
	if err != nil {
		return nil, err
	}

	resOut, err := singleOutport(runtimeIO, "res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			msg, received := dataIn.Receive(ctx)
			if !received {
				return
			}

			if !forwardForEachMessage(ctx, doneIn, itemOut, resOut, msg) {
				return
			}
		}
	}, nil
}

func forwardForEachMessage(
	ctx context.Context,
	doneIn runtime.SingleInport,
	itemOut, resOut runtime.SingleOutport,
	msg runtime.Msg,
) bool {
	switch {
	case isStreamOpen(msg), isStreamClose(msg):
		return resOut.Send(ctx, msg)
	case isStreamData(msg):
		if !itemOut.Send(ctx, streamDataValue(msg)) {
			return false
		}

		if _, received := doneIn.Receive(ctx); !received {
			return false
		}

		return resOut.Send(ctx, msg)
	default:
		panic("stream_for_each_controller: unexpected stream tag")
	}
}
