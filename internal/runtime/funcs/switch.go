package funcs

import (
	"context"
	"errors"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type switchRouter struct{}

//nolint:cyclop,gocognit,gocyclo,varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (switchRouter) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	caseArrIn, err := io.In.Array("case")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	caseOut, err := io.Out.Array("case")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	elseOut, err := io.Out.Single("else")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	if caseArrIn.Len() != caseOut.Len() {
		return nil, errors.New("number of 'case' inports must match number of outports")
	}

	return func(ctx context.Context) {
		for {
			var (
				//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
				wg              sync.WaitGroup
				dataMsg         runtime.Msg
				cases           = make([]runtime.Msg, caseArrIn.Len())
				dataOk, casesOk bool
			)

			wg.Go(func() {
				dataMsg, dataOk = dataIn.Receive(ctx)
			})

			wg.Go(func() {
				casesOk = caseArrIn.ReceiveAll(ctx, func(idx int, msg runtime.Msg) bool {
					cases[idx] = msg
					return true
				})
			})

			wg.Wait()

			if !dataOk || !casesOk {
				return
			}

			matchIdx := -1
			for i, caseMsg := range cases {
				if runtime.Match(dataMsg, caseMsg) {
					matchIdx = i
					break
				}
			}

			if matchIdx != -1 {
				caseIdx := runtime.Uint8Index(matchIdx)
				if !caseOut.Send(
					ctx,
					caseIdx,
					tryToUnboxIfUnion(dataMsg),
				) {
					return
				}
				continue
			}

			// For unions: we never unbox even if possible when sending to :else
			if !elseOut.Send(ctx, dataMsg) {
				return
			}
		}
	}, nil
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func tryToUnboxIfUnion(dataMsg runtime.Msg) runtime.Msg {
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	u, ok := dataMsg.(runtime.UnionMsg)
	if !ok {
		return dataMsg
	}

	if u.Data() == nil {
		return u
	}

	return u.Data()
}
