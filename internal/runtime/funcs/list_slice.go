package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type listSlice struct{}

func (listSlice) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

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
			var dataMsg, fromMsg, toMsg runtime.Msg
			var dataOK, fromOK, toOK bool

			var wg sync.WaitGroup
			wg.Go(func() {
				dataMsg, dataOK = dataIn.Receive(ctx)
			})
			wg.Go(func() {
				fromMsg, fromOK = fromIn.Receive(ctx)
			})
			wg.Go(func() {
				toMsg, toOK = toIn.Receive(ctx)
			})
			wg.Wait()

			if !dataOK || !fromOK || !toOK {
				return
			}

			data := dataMsg.List()
			l := int64(len(data))

			from := fromMsg.Int()
			if from < 0 {
				from += l
			}
			to := toMsg.Int()
			if to < 0 {
				to += l
			}

			if from < 0 || to < 0 || from > to || to > l {
				panic("slice index out of bounds")
			}

			res := runtime.NewListMsg(append([]runtime.Msg(nil), data[from:to]...))
			if !resOut.Send(ctx, res) {
				return
			}
		}
	}, nil
}
