package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type rangeInt struct{}

func (rangeInt) Create(runtimeIO runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	fromIn, err := singleInport(runtimeIO, "from")
	if err != nil {
		return nil, err
	}

	toIn, err := singleInport(runtimeIO, "to")
	if err != nil {
		return nil, err
	}

	resOut, err := singleOutport(runtimeIO, "res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			fromMsg, toMsg, received := receive2(ctx, fromIn, toIn)
			if !received {
				return
			}

			from := fromMsg.Int()
			toValue := toMsg.Int()

			if !sendIntRange(ctx, resOut, from, toValue) {
				return
			}
		}
	}, nil
}

func sendIntRange(
	ctx context.Context,
	resOut runtime.SingleOutport,
	from, toValue int64,
) bool {
	if !resOut.Send(ctx, streamOpen()) {
		return false
	}

	if from < toValue {
		return sendAscendingIntRange(ctx, resOut, from, toValue)
	}

	if from > toValue {
		return sendDescendingIntRange(ctx, resOut, from, toValue)
	}

	return resOut.Send(ctx, streamClose())
}

func sendAscendingIntRange(
	ctx context.Context,
	resOut runtime.SingleOutport,
	from, toValue int64,
) bool {
	for data := from; data < toValue; data++ {
		if !resOut.Send(ctx, streamData(runtime.NewIntMsg(data))) {
			return false
		}
	}

	return resOut.Send(ctx, streamClose())
}

func sendDescendingIntRange(
	ctx context.Context,
	resOut runtime.SingleOutport,
	from, toValue int64,
) bool {
	for data := from; data > toValue; data-- {
		if !resOut.Send(ctx, streamData(runtime.NewIntMsg(data))) {
			return false
		}
	}

	return resOut.Send(ctx, streamClose())
}
