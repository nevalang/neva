package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

// streamMapController preserves stream event ordering for Map.
// For each Data event it waits for mapped payload before forwarding Data.
type streamMapController struct{}

func (streamMapController) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	mappedIn, err := io.In.Single("mapped")
	if err != nil {
		return nil, err
	}

	itemOut, err := io.Out.Single("item")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			msg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			switch {
			case isStreamOpen(msg):
				if !resOut.Send(ctx, streamOpen()) {
					return
				}
			case isStreamData(msg):
				if !itemOut.Send(ctx, streamDataValue(msg)) {
					return
				}

				mappedMsg, ok := mappedIn.Receive(ctx)
				if !ok {
					return
				}

				if !resOut.Send(ctx, streamData(mappedMsg)) {
					return
				}
			case isStreamClose(msg):
				if !resOut.Send(ctx, streamClose()) {
					return
				}
			default:
				panic("stream_map_controller: unexpected stream tag")
			}
		}
	}, nil
}
