package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type listSlice struct{}

// sliceList returns a copy of a normalized list slice.
func sliceList(data []runtime.Msg, from int64, to int64) []runtime.Msg {
	start, end := normalizeSliceBounds(from, to, int64(len(data)))
	return append([]runtime.Msg(nil), data[start:end]...)
}

//nolint:dupl,varnamelen
func (listSlice) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	fromIn, err := io.In.Single("from")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	toIn, err := io.In.Single("to")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var dataMsg, fromMsg, toMsg runtime.Msg
			var dataOK, fromOK, toOK bool

			//nolint:varnamelen
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

			res := runtime.NewListMsg(sliceList(dataMsg.List(), fromMsg.Int(), toMsg.Int()))
			if !resOut.Send(ctx, res) {
				return
			}
		}
	}, nil
}
