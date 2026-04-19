//nolint:dupl // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type stringAdd struct{}

func (stringAdd) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	leftIn, err := singleIn(io, "left")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	rightIn, err := singleIn(io, "right")
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
			var leftMsg, rightMsg runtime.Msg
			var leftOk, rightOk bool
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			var wg sync.WaitGroup

			wg.Go(func() {
				leftMsg, leftOk = leftIn.Receive(ctx)
			})

			wg.Go(func() {
				rightMsg, rightOk = rightIn.Receive(ctx)
			})

			wg.Wait()

			if !leftOk || !rightOk {
				return
			}

			resMsg := runtime.NewStringMsg(leftMsg.Str() + rightMsg.Str())
			if !resOut.Send(ctx, resMsg) {
				return
			}
		}
	}, nil
}
