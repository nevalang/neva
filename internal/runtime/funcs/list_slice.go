package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type listSlice struct{}

// sliceList returns a copy of a normalized list slice.
func sliceList(data []messages.Msg, from int64, to int64) []messages.Msg {
	start, end := normalizeSliceBounds(from, to, int64(len(data)))
	return append([]messages.Msg(nil), data[start:end]...)
}

func sliceTypedList[T any](data []T, from int64, to int64) []T {
	start, end := normalizeSliceBounds(from, to, int64(len(data)))
	return append([]T(nil), data[start:end]...)
}

func (listSlice) Create(io runtime.IO, _ messages.Msg) (func(context.Context), error) {
	dataIn, fromIn, toIn, resOut, err := resolveListSlicePorts(io)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, fromMsg, toMsg, received := receive3(ctx, dataIn, fromIn, toIn)
			if !received {
				return
			}

			if !resOut.Send(ctx, sliceMessage(dataMsg.Msg, fromMsg.Int(), toMsg.Int())) {
				return
			}
		}
	}, nil
}

func resolveListSlicePorts(
	runtimeIO runtime.IO,
) (
	runtime.SingleInport,
	runtime.SingleInport,
	runtime.SingleInport,
	runtime.SingleOutport,
	error,
) {
	dataIn, err := getSingleInport(runtimeIO.In, "data")
	if err != nil {
		return runtime.SingleInport{}, runtime.SingleInport{}, runtime.SingleInport{}, runtime.SingleOutport{}, err
	}

	fromIn, err := getSingleInport(runtimeIO.In, "from")
	if err != nil {
		return runtime.SingleInport{}, runtime.SingleInport{}, runtime.SingleInport{}, runtime.SingleOutport{}, err
	}

	toIn, err := getSingleInport(runtimeIO.In, "to")
	if err != nil {
		return runtime.SingleInport{}, runtime.SingleInport{}, runtime.SingleInport{}, runtime.SingleOutport{}, err
	}

	resOut, err := getSingleOutport(runtimeIO.Out, "res")
	if err != nil {
		return runtime.SingleInport{}, runtime.SingleInport{}, runtime.SingleInport{}, runtime.SingleOutport{}, err
	}

	return dataIn, fromIn, toIn, resOut, nil
}

//nolint:ireturn // Runtime messages have multiple concrete representations.
func sliceMessage(data messages.Msg, start int64, end int64) messages.Msg {
	list := data.List()
	if values, ok := messages.AsListInts(list); ok {
		return messages.NewListIntMsg(sliceTypedList(values, start, end))
	}
	if values, ok := messages.AsListStrings(list); ok {
		return messages.NewListStringMsg(sliceTypedList(values, start, end))
	}
	if values, ok := messages.AsListBools(list); ok {
		return messages.NewListBoolMsg(sliceTypedList(values, start, end))
	}
	if values, ok := messages.AsListFloats(list); ok {
		return messages.NewListFloatMsg(sliceTypedList(values, start, end))
	}

	return messages.NewListMsg(sliceList(list.Untyped(), start, end))
}
