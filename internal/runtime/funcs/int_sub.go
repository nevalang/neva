package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intSub struct{}

func (intSub) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	seqIn, err := io.In.Single("seq")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.SingleOutport("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var (
			acc     int64 = 0
			started bool  = false
		)

		for {
			seqMsg, ok := seqIn.Receive(ctx)
			if !ok {
				return
			}

			item := seqMsg.Map()

			if !started {
				acc = item["data"].Int()
				started = true
			} else {
				acc -= item["data"].Int()
			}

			if item["last"].Bool() {
				if !resOut.Send(ctx, runtime.NewIntMsg(acc)) {
					return
				}
				acc = 0
				started = false
			}
		}
	}, nil
}
