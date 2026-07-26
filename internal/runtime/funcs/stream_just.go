package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type streamJust struct{}

func (streamJust) Create(runtimeIO runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	dataIn, err := singleInport(runtimeIO, "data")
	if err != nil {
		return nil, err
	}

	resOut, err := singleOutport(runtimeIO, "res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, received := dataIn.Receive(ctx)
			if !received {
				return
			}

			if !sendSingleItemStream(ctx, resOut, dataMsg) {
				return
			}
		}
	}, nil
}

func sendSingleItemStream(ctx context.Context, resOut runtime.SingleOutport, dataMsg runtime.OrderedMsg) bool {
	return resOut.Send(ctx, messages.NewStreamOpenMsg()) &&
		resOut.Send(ctx, messages.NewStreamDataMsg(dataMsg.Msg)) &&
		resOut.Send(ctx, messages.NewStreamCloseMsg())
}
