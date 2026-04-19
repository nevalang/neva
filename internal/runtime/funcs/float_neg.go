package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type floatNeg struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (floatNeg) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
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

			if !resOut.Send(ctx, runtime.NewFloatMsg(-dataMsg.Float())) {
				return
			}
		}
	}, nil
}
