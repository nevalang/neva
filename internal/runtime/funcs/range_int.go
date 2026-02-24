package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type rangeInt struct{}

func (rangeInt) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	fromIn, err := io.In.Single("from")
	if err != nil {
		return nil, err
	}

	toIn, err := io.In.Single("to")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			fromMsg, ok := fromIn.Receive(ctx)
			if !ok {
				return
			}

			toMsg, ok := toIn.Receive(ctx)
			if !ok {
				return
			}

			from := fromMsg.Int()
			to := toMsg.Int()

			if !resOut.Send(ctx, streamOpen()) {
				return
			}

			if from < to {
				for data := from; data < to; data++ {
					if !resOut.Send(ctx, streamData(runtime.NewIntMsg(data))) {
						return
					}
				}
			} else if from > to {
				for data := from; data > to; data-- {
					if !resOut.Send(ctx, streamData(runtime.NewIntMsg(data))) {
						return
					}
				}
			}

			if !resOut.Send(ctx, streamClose()) {
				return
			}
		}
	}, nil
}
