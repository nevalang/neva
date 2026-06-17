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

func sliceTypedList[T any](data []T, from int64, to int64) []T {
	start, end := normalizeSliceBounds(from, to, int64(len(data)))
	return append([]T(nil), data[start:end]...)
}

//nolint:dupl,varnamelen,gocognit // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (listSlice) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	fromIn, err := io.In.Single("from")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	toIn, err := io.In.Single("to")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var dataMsg, fromMsg, toMsg runtime.Msg
			var dataOK, fromOK, toOK bool

			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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

			list := dataMsg.List()
			from, to := fromMsg.Int(), toMsg.Int()

			var sent bool
			//nolint:nestif // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			if values, ok := runtime.AsListInts(list); ok {
				sent = resOut.Send(ctx, runtime.NewListIntMsg(sliceTypedList(values, from, to)))
			} else if values, ok := runtime.AsListStrings(list); ok {
				sent = resOut.Send(ctx, runtime.NewListStringMsg(sliceTypedList(values, from, to)))
			} else if values, ok := runtime.AsListBools(list); ok {
				sent = resOut.Send(ctx, runtime.NewListBoolMsg(sliceTypedList(values, from, to)))
			} else if values, ok := runtime.AsListFloats(list); ok {
				sent = resOut.Send(ctx, runtime.NewListFloatMsg(sliceTypedList(values, from, to)))
			} else {
				sent = resOut.Send(ctx, runtime.NewListMsg(sliceList(list.Msgs(), from, to)))
			}

			if !sent {
				return
			}
		}
	}, nil
}
