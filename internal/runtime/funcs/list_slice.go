package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type listSlice struct{}

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

			if !resOut.Send(ctx, messages.ListSlice(dataMsg.List(), fromMsg.Int(), toMsg.Int())) {
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
