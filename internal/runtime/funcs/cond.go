package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type cond struct{}

//nolint:gocognit,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (c cond) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := singleIn(io, "data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	ifIn, err := singleIn(io, "if")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	thenOut, err := singleOut(io, "then")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	elseOut, err := singleOut(io, "else")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var dataMsg, ifMsg runtime.Msg
			var dataOk, ifOk bool

			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			var wg sync.WaitGroup

			wg.Go(func() {
				dataMsg, dataOk = dataIn.Receive(ctx)
			})

			wg.Go(func() {
				ifMsg, ifOk = ifIn.Receive(ctx)
			})

			wg.Wait()

			if !dataOk || !ifOk {
				return
			}

			var out runtime.SingleOutport
			if ifMsg.Bool() {
				out = thenOut
			} else {
				out = elseOut
			}

			if !out.Send(ctx, dataMsg) {
				return
			}
		}
	}, nil
}
