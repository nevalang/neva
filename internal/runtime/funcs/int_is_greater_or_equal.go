//nolint:dupl // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type intIsGreaterOrEqual struct{}

func (intIsGreaterOrEqual) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	accIn, err := singleIn(io, "left")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	elIn, err := singleIn(io, "right")
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
			var accMsg, elMsg runtime.Msg
			var accOk, elOk bool

			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			var wg sync.WaitGroup

			wg.Go(func() {
				accMsg, accOk = accIn.Receive(ctx)
			})

			wg.Go(func() {
				elMsg, elOk = elIn.Receive(ctx)
			})

			wg.Wait()

			if !accOk || !elOk {
				return
			}

			if !resOut.Send(ctx, runtime.NewBoolMsg(accMsg.Int() >= elMsg.Int())) {
				return
			}
		}
	}, nil
}
