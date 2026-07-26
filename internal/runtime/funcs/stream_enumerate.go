package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type streamEnumerate struct{}

func (streamEnumerate) Create(runtimeIO runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
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

			if !forwardEnumeratedMessage(ctx, resOut, msg.Msg, &idx) {
				return
			}
		}
	}, nil
}

func forwardEnumeratedMessage(
	ctx context.Context,
	resOut runtime.SingleOutport,
	msg messages.Msg,
	idx *int64,
) bool {
	switch {
	case messages.IsStreamOpen(msg):
		*idx = 0
		return resOut.Send(ctx, messages.NewStreamOpenMsg())
	case messages.IsStreamData(msg):
		// Enumerated<T> is the Data union payload, so encode it as a struct message first.
		item := messages.NewStructMsg([]messages.StructField{
			messages.NewStructField("idx", messages.NewIntMsg(*idx)),
			messages.NewStructField("item", messages.StreamDataValue(msg)),
		})
		if !resOut.Send(ctx, messages.NewStreamDataMsg(item)) {
			return false
		}

		*idx++
		return true
	case messages.IsStreamClose(msg):
		return resOut.Send(ctx, messages.NewStreamCloseMsg())
	default:
		panic("stream_enumerate: unexpected stream tag")
	}
}
