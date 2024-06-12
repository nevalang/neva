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
	fromIn, err := io.In.SingleInport("from")
	if err != nil {
		return nil, err
	}

	toIn, err := io.In.SingleInport("to")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.SingleOutport("data")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var fromMsg, toMsg runtime.Msg

		for {
			select {
			case <-ctx.Done():
				return
			case fromMsg = <-fromIn:
			}

			select {
			case <-ctx.Done():
				return
			case toMsg = <-toIn:
			}
			var (
				idx  int64 = 0
				last bool  = false
				data int64 = fromMsg.Int()
			)

			if fromMsg.Int() < toMsg.Int() {
				for !last {
					if data == toMsg.Int()-1 {
						last = true
					}

					item := streamItem(
						runtime.NewIntMsg(data),
						idx,
						last,
					)

					select {
					case <-ctx.Done():
						return
					case dataOut <- item:
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

					select {
					case <-ctx.Done():
						return
					case dataOut <- item:
					}

					idx++
					data--
				}
			}

		}
	}, nil
}
