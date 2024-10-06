package funcs

import (
	"context"
	"errors"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type switcher struct{}

func (switcher) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	caseIn, err := io.In.Array("case")
	if err != nil {
		return nil, err
	}

	caseOut, err := io.Out.Array("case")
	if err != nil {
		return nil, err
	}

	elseOut, err := io.Out.Single("else")
	if err != nil {
		return nil, err
	}

	if caseIn.Len() != caseOut.Len() {
		return nil, errors.New("number of 'case' inports must match number of outports")
	}

	return func(ctx context.Context) {
		for {
			var (
				wg              sync.WaitGroup
				dataMsg         runtime.Msg
				cases           = make([]runtime.Msg, caseIn.Len())
				dataOk, casesOk bool
			)

			wg.Add(2)

			go func() {
				dataMsg, dataOk = dataIn.Receive(ctx)
				wg.Done()
			}()

			go func() {
				casesOk = caseIn.ReceiveAll(ctx, func(idx int, msg runtime.Msg) bool {
					cases[idx] = msg
					return true
				})
				wg.Done()
			}()

			wg.Wait()

			if !dataOk || !casesOk {
				return
			}

			matchIdx := -1
			for i, caseMsg := range cases {
				if dataMsg.Equal(caseMsg) {
					matchIdx = i
					break
				}
			}

			if matchIdx != -1 {
				if !caseOut.Send(ctx, uint8(matchIdx), dataMsg) {
					return
				}
				continue
			}

			if !elseOut.Send(ctx, dataMsg) {
				return
			}
		}
	}, nil
}
