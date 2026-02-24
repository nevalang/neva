package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamEnumerate struct{}

func (streamEnumerate) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var idx int64
		for {
			msg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			switch {
			case isStreamOpen(msg):
				idx = 0
				if !resOut.Send(ctx, streamOpen()) {
					return
				}
			case isStreamData(msg):
				item := runtime.NewStructMsg([]runtime.StructField{
					runtime.NewStructField("idx", runtime.NewIntMsg(idx)),
					runtime.NewStructField("data", streamDataValue(msg)),
				})
				if !resOut.Send(ctx, streamData(item)) {
					return
				}
				idx++
			case isStreamClose(msg):
				if !resOut.Send(ctx, streamClose()) {
					return
				}
			}
		}
	}, nil
}
