//nolint:dupl // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
package funcs

import (
	"context"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type floatDiv struct{}

func (floatDiv) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	accIn, err := io.In.Single("left")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	elIn, err := io.In.Single("right")
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

			resMsg := runtime.NewFloatMsg(accMsg.Float() / elMsg.Float())
			if !resOut.Send(ctx, resMsg) {
				return
			}
		}
	}, nil
}
