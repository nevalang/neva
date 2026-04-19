package funcs

import (
	"context"
	"errors"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type race struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (race) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := singleIn(io, "data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	casesArrIn, err := arrayIn(io, "case")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	casesOut, err := arrayOut(io, "case")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	if casesArrIn.Len() != casesOut.Len() {
		return nil, errors.New("number of 'case' inports must match number of 'case' outports")
	}

	return func(ctx context.Context) {
		var (
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			wg      sync.WaitGroup
			dataMsg runtime.Msg
			dataOk  bool
			caseMsg runtime.SelectedMsg
			caseOk  bool
		)
		for {
			wg.Go(func() {
				dataMsg, dataOk = dataIn.Receive(ctx)
			})
			wg.Go(func() {
				caseMsg, caseOk = casesArrIn.Select(ctx)
			})
			wg.Wait()
			if !dataOk || !caseOk {
				return
			}
			if !casesOut.Send(ctx, caseMsg.SlotIdx, dataMsg) {
				return
			}
		}
	}, nil
}
