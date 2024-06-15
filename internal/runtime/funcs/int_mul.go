package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intMul struct{}

func (intMul) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	seqIn, err := io.In.Single("seq")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.SingleOutport("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var acc int64 = 1

		for {
			seqMsg, ok := seqIn.Receive(ctx)
			if !ok {
				return
			}

			item := seqMsg.Map()
			acc *= item["data"].Int()

			if !item["last"].Bool() {
				continue
			}

			if !resOut.Send(ctx, runtime.NewIntMsg(acc)) {
				return
			}

			acc = 1
		}
	}, nil
}
