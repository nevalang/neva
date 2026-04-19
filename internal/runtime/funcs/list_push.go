package funcs

import (
	"context"
	"slices"

	"github.com/nevalang/neva/internal/runtime"
)

type listPush struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (p listPush) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := singleIn(io, "data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}
	lstIn, err := singleIn(io, "lst")
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
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			lstMsg, ok := lstIn.Receive(ctx)
			if !ok {
				return
			}

			lstCopy := slices.Clone(lstMsg.List())

			if !resOut.Send(
				ctx,
				runtime.NewListMsg(
					append(lstCopy, dataMsg),
				),
			) {
				return
			}
		}
	}, nil
}
