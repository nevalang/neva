package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intAdd struct{}

func (intAdd) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	seqIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var acc int64 = 0

		for {
			msg, ok := seqIn.Receive(ctx)
			if !ok {
				return
			}

			item := msg.Map()

			acc += item["data"].Int()

			if !item["last"].Bool() {
				continue
			}

			if !resOut.Send(ctx, runtime.NewIntMsg(acc)) {
				return
			}

			acc = 0
		}
	}, nil
}
