package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intDec struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (i intDec) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	dataIn, err := singleIn(io, "data")
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
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(ctx, runtime.NewIntMsg(dataMsg.Int()-1)) {
				return
			}
		}
	}, nil
}
