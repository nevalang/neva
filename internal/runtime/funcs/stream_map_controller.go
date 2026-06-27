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
	case runtime.IsStreamOpen(msg):
		return resOut.Send(ctx, runtime.NewStreamOpenMsg())
	case runtime.IsStreamData(msg):
		if !itemOut.Send(ctx, runtime.StreamDataValue(msg)) {
			return false
		}

		mappedMsg, received := mappedIn.Receive(ctx)
		if !received {
			return false
		}

		return resOut.Send(ctx, runtime.NewStreamDataMsg(mappedMsg.Msg))
	case runtime.IsStreamClose(msg):
		return resOut.Send(ctx, runtime.NewStreamCloseMsg())
	default:
		panic("stream_map_controller: unexpected stream tag")
	}
}
