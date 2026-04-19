package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type lock struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (l lock) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := io.In.Single("sig")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	dataIn, err := io.In.Single("data")
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
			var (
				//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
				wg            sync.WaitGroup
				data          runtime.Msg
				dataOk, sigOk bool
			)

			wg.Go(func() {
				data, dataOk = dataIn.Receive(ctx)
			})
			wg.Go(func() {
				_, sigOk = sigIn.Receive(ctx)
			})
			wg.Wait()

			if !dataOk || !sigOk {
				return
			}

			if !resOut.Send(ctx, data) {
				return
			}
		}
	}, nil
}
