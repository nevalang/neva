package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamZipMany struct{}

func (streamZipMany) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	dataIn, err := io.In.Array("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		streamsCount := dataIn.Len()
		index := int64(0)

		for {
			zipped := make([]runtime.Msg, streamsCount)
			shouldStop := false

			for streamIdx := 0; streamIdx < streamsCount; streamIdx++ {
				msg, ok := dataIn.Receive(ctx, streamIdx)
				if !ok {
					return
				}

				item := msg.Struct()
				zipped[streamIdx] = item.Get("data")

				if item.Get("last").Bool() {
					shouldStop = true
				}
			}

			if !resOut.Send(
				ctx,
				streamItem(
					runtime.NewListMsg(zipped),
					index,
					shouldStop,
				),
			) {
				return
			}

			if shouldStop {
				return
			}

			index++
		}
	}, nil
}
