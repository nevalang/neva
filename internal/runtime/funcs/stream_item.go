package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamItem struct{}

func (streamItem) Create(runtimeIO runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
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
	return resOut.Send(ctx, newStreamOpenMsg()) &&
		resOut.Send(ctx, newStreamDataMsg(dataMsg.Msg)) &&
		resOut.Send(ctx, newStreamCloseMsg())
}
