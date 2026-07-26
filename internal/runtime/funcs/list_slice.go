package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type listSlice struct{}

// sliceList returns a copy of a normalized list slice.
func sliceList(data []messages.Msg, from int64, to int64) []messages.Msg {
	start, end := normalizeSliceBounds(from, to, int64(len(data)))
	return append([]messages.Msg(nil), data[start:end]...)
}

func sliceTypedList[T any](data []T, from int64, to int64) []T {
	start, end := normalizeSliceBounds(from, to, int64(len(data)))
	return append([]T(nil), data[start:end]...)
}

//nolint:dupl,varnamelen,gocognit // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (listSlice) Create(io runtime.IO, _ messages.Msg) (func(context.Context), error) {
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
			var dataMsg, fromMsg, toMsg messages.Msg
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
			if values, ok := messages.AsListInts(list); ok {
				sent = resOut.Send(ctx, messages.NewListIntMsg(sliceTypedList(values, from, to)))
			} else if values, ok := messages.AsListStrings(list); ok {
				sent = resOut.Send(ctx, messages.NewListStringMsg(sliceTypedList(values, from, to)))
			} else if values, ok := messages.AsListBools(list); ok {
				sent = resOut.Send(ctx, messages.NewListBoolMsg(sliceTypedList(values, from, to)))
			} else if values, ok := messages.AsListFloats(list); ok {
				sent = resOut.Send(ctx, messages.NewListFloatMsg(sliceTypedList(values, from, to)))
			} else {
				sent = resOut.Send(ctx, messages.NewListMsg(sliceList(list.Untyped(), from, to)))
			}

			if !sent {
				return
			}
		}
	}, nil
}
