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
			fromMsg, fromReceived := fromIn.Receive(ctx)
			if !fromReceived {
				return
			}

			toMsg, toReceived := toIn.Receive(ctx)
			if !toReceived {
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

	for _, data := range rangeValues(from, toValue) {
		if !resOut.Send(ctx, streamData(runtime.NewIntMsg(data))) {
			return false
		}
	}

	return resOut.Send(ctx, streamClose())
}

func rangeValues(from, toValue int64) []int64 {
	switch {
	case from < toValue:
		values := make([]int64, 0, toValue-from)
		for data := from; data < toValue; data++ {
			values = append(values, data)
		}

		return values
	case from > toValue:
		values := make([]int64, 0, from-toValue)
		for data := from; data > toValue; data-- {
			values = append(values, data)
		}

		return values
	default:
		return nil
	}
}
