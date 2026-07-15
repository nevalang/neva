package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamSplitController struct{}

func (streamSplitController) Create(runtimeIO runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := singleInport(runtimeIO, "data")
	if err != nil {
		return nil, err
	}

	passIn, err := singleInport(runtimeIO, "pass")
	if err != nil {
		return nil, err
	}

	itemOut, err := singleOutport(runtimeIO, "item")
	if err != nil {
		return nil, err
	}

	thenOut, err := singleOutport(runtimeIO, "then")
	if err != nil {
		return nil, err
	}

	elseOut, err := singleOutport(runtimeIO, "else")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			msg, received := dataIn.Receive(ctx)
			if !received {
				return
			}

			if !forwardSplitMessage(ctx, passIn, itemOut, thenOut, elseOut, msg) {
				return
			}
		}
	}, nil
}

func forwardSplitMessage(
	ctx context.Context,
	passIn runtime.SingleInport,
	itemOut runtime.SingleOutport,
	thenOut, elseOut runtime.SingleOutport,
	msg runtime.Msg,
) bool {
	switch {
	case isStreamOpen(msg):
		return thenOut.Send(ctx, newStreamOpenMsg()) && elseOut.Send(ctx, newStreamOpenMsg())
	case isStreamData(msg):
		item := streamDataValue(msg)
		if !itemOut.Send(ctx, item) {
			return false
		}

		passMsg, received := passIn.Receive(ctx)
		if !received {
			return false
		}

		if passMsg.Bool() {
			return thenOut.Send(ctx, msg)
		}

		return elseOut.Send(ctx, msg)
	case isStreamClose(msg):
		return thenOut.Send(ctx, newStreamCloseMsg()) && elseOut.Send(ctx, newStreamCloseMsg())
	default:
		panic("stream_split_controller: unexpected stream tag")
	}
}
