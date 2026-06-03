package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

// streamMapController preserves stream event ordering for Map.
// For each Data event it waits for mapped payload before forwarding Data.
type streamMapController struct{}

func (streamMapController) Create(runtimeIO runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := singleInport(runtimeIO, "data")
	if err != nil {
		return nil, err
	}

	mappedIn, err := singleInport(runtimeIO, "mapped")
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

			if !forwardMappedMessage(ctx, mappedIn, itemOut, resOut, msg) {
				return
			}
		}
	}, nil
}

func forwardMappedMessage(
	ctx context.Context,
	mappedIn runtime.SingleInport,
	itemOut, resOut runtime.SingleOutport,
	msg runtime.Msg,
) bool {
	switch {
	case isStreamOpen(msg):
		return resOut.Send(ctx, streamOpen())
	case isStreamData(msg):
		if !itemOut.Send(ctx, streamDataValue(msg)) {
			return false
		}

		mappedMsg, received := mappedIn.Receive(ctx)
		if !received {
			return false
		}

		return resOut.Send(ctx, streamData(mappedMsg.Msg))
	case isStreamClose(msg):
		return resOut.Send(ctx, streamClose())
	default:
		panic("stream_map_controller: unexpected stream tag")
	}
}
