//nolint:dupl // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type and struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (p and) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	aIn, err := singleIn(io, "left")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	bIn, err := singleIn(io, "right")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := singleOut(io, "res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	// TODO send false as soon as A in is false, but do it correctly
	return func(ctx context.Context) {
		for {
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			aMsg, ok := aIn.Receive(ctx)
			if !ok {
				return
			}

			bMsg, ok := bIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(
				ctx,
				runtime.NewBoolMsg(aMsg.Bool() && bMsg.Bool()),
			) {
				return
			}
		}
	}, nil
}
