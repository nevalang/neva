package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type notEq struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (p notEq) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	actualIn, err := singleIn(io, "left")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	comparedIn, err := singleIn(io, "right")
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
			var (
				//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
				wg                   sync.WaitGroup
				val1, val2           runtime.Msg
				actualOk, comparedOk bool
			)

			wg.Go(func() {
				val1, actualOk = actualIn.Receive(ctx)
			})
			wg.Go(func() {
				val2, comparedOk = comparedIn.Receive(ctx)
			})
			wg.Wait()

			if !actualOk || !comparedOk {
				return
			}

			if !resOut.Send(
				ctx,
				runtime.NewBoolMsg(!val1.Equal(val2)),
			) {
				return
			}
		}
	}, nil
}
