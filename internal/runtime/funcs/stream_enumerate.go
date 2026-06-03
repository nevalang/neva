package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamEnumerate struct{}

func (streamEnumerate) Create(runtimeIO runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := singleInport(runtimeIO, "data")
	if err != nil {
		return nil, err
	}

	resOut, err := singleOutport(runtimeIO, "res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var idx int64
		for {
			msg, received := dataIn.Receive(ctx)
			if !received {
				return
			}

			if !forwardEnumeratedMessage(ctx, resOut, msg, &idx) {
				return
			}
		}
	}, nil
}

func forwardEnumeratedMessage(
	ctx context.Context,
	resOut runtime.SingleOutport,
	msg runtime.Msg,
	idx *int64,
) bool {
	switch {
	case isStreamOpen(msg):
		*idx = 0
		return resOut.Send(ctx, streamOpen())
	case isStreamData(msg):
		item := runtime.NewStructMsg([]runtime.StructField{
			runtime.NewStructField("idx", runtime.NewIntMsg(*idx)),
			runtime.NewStructField("item", streamDataValue(msg)),
		})
		if !resOut.Send(ctx, streamData(item)) {
			return false
		}

		*idx++
		return true
	case isStreamClose(msg):
		return resOut.Send(ctx, streamClose())
	default:
		panic("stream_enumerate: unexpected stream tag")
	}
}
