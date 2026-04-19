package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type stringSlice struct{}

// sliceString returns a normalized slice over rune indices.
func sliceString(data string, from int64, to int64) string {
	runes := []rune(data)
	start, end := normalizeSliceBounds(from, to, int64(len(runes)))
	return string(runes[start:end])
}

//nolint:dupl,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (stringSlice) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	dataIn, err := singleIn(io, "data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	fromIn, err := singleIn(io, "from")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	toIn, err := singleIn(io, "to")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := singleOut(io, "res")
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

			res := runtime.NewStringMsg(sliceString(dataMsg.Str(), fromMsg.Int(), toMsg.Int()))
			if !resOut.Send(ctx, res) {
				return
			}
		}
	}, nil
}
