package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type unionWrapper struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (unionWrapper) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := singleIn(io, "data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	tagIn, err := singleIn(io, "tag")
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

			tagMsg, ok := tagIn.Receive(ctx)
			if !ok {
				return
			}

			tag := tagMsg.Union().Tag()
			if !resOut.Send(ctx, runtime.NewUnionMsg(tag, dataMsg)) {
				return
			}
		}
	}, nil
}
