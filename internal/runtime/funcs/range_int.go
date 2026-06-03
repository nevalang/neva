package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type rangeInt struct{}

//nolint:cyclop,gocognit,gocyclo,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (rangeInt) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
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
			fromMsg, toMsg, ok := receive2(ctx, fromIn, toIn)
			if !ok {
				return
			}

			var (
				from = fromMsg.Int()
				//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
				to  = toMsg.Int()
				idx = int64(0)
			)

			if from < to {
				for data := from; data < to; data++ {
					item := streamItem(
						runtime.NewIntMsg(data),
						idx,
						data == to-1,
					)

					if !resOut.Send(ctx, item) {
						return
					}

					idx++
				}
			} else {
				for data := from; data > to; data-- {
					item := streamItem(
						runtime.NewIntMsg(data),
						idx,
						data == to+1,
					)

					if !resOut.Send(ctx, item) {
						return
					}

					idx++
				}
			}
		}
	}, nil
}
