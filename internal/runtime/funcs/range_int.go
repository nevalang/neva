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
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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
				//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
				to = toMsg.Int()

				idx  = int64(0)
				last = false
				data = from
			)

			//nolint:nestif // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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

					if !resOut.Send(ctx, item) {
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

					if !resOut.Send(ctx, item) {
						return
					}

					idx++
					data--
				}
			}
		}
	}, nil
}
