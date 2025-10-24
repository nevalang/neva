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
	streamsIn, err := io.In.Array("streams")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Single("data")
	if err != nil {
		return nil, err
	}

	// If there are no streams connected there is nothing to zip.
	if streamsIn.Len() == 0 {
		return func(ctx context.Context) {}, nil
	}

	return func(ctx context.Context) {
		streamsCount := streamsIn.Len()
		index := int64(0)

		for {
			zipped := make([]runtime.Msg, streamsCount)
			shouldStop := false

			for streamIdx := 0; streamIdx < streamsCount; streamIdx++ {
				msg, ok := streamsIn.Receive(ctx, streamIdx)
				if !ok {
					return
				}

				item := msg.Struct()
				zipped[streamIdx] = item.Get("data")

				if item.Get("last").Bool() {
					shouldStop = true
				}
			}

			if !dataOut.Send(
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
