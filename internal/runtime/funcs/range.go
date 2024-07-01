package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamIntRange struct{}

func (streamIntRange) Create(
	io runtime.FuncIO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	fromIn, err := io.In.Single("from")
	if err != nil {
		return nil, err
	}

	toIn, err := io.In.Single("to")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Single("data")
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

			var (
				from = fromMsg.Int()
				to   = toMsg.Int()

				idx  int64 = 0
				last bool  = false
				data int64 = from
			)

			if from < to {
				for !last {
					if data == to-1 {
						last = true
					}

					item := streamItem(
						runtime.NewIntMsg(data),
						idx,
						last,
					)

					if !dataOut.Send(ctx, item) {
						return
					}

					idx++
					data++
				}
			} else {
				for !last {
					if data == toMsg.Int()+1 {
						last = true
					}

					item := streamItem(
						runtime.NewIntMsg(data),
						idx,
						last,
					)

					if !dataOut.Send(ctx, item) {
						return
					}

					idx++
					data--
				}
			}
		}
	}, nil
}
