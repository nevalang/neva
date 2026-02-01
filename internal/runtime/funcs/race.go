package funcs

import (
	"context"
	"errors"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type race struct{}

func (race) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	casesArrIn, err := io.In.Array("case")
	if err != nil {
		return nil, err
	}

	casesOut, err := io.Out.Array("case")
	if err != nil {
		return nil, err
	}

	if casesArrIn.Len() != casesOut.Len() {
		return nil, errors.New("number of 'case' inports must match number of 'case' outports")
	}

	return func(ctx context.Context) {
		var (
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
