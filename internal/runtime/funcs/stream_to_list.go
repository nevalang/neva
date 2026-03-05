package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamToList struct{}

func (s streamToList) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	seqIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		// Fully materializes one stream batch before emitting resulting list.
		list := []runtime.Msg{}

		for {
			msg, ok := seqIn.Receive(ctx)
			if !ok {
				return
			}

			switch {
			case isStreamOpen(msg):
				list = list[:0]
				continue
			case isStreamData(msg):
				list = append(list, streamDataValue(msg))
				continue
			case !isStreamClose(msg):
				continue
			}

			if !resOut.Send(ctx, runtime.NewListMsg(list)) {
				return
			}
		}
	}, nil
}
