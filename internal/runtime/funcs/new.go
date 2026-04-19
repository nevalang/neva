package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type newV2 struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (c newV2) Create(io runtime.IO, cfg runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := singleIn(io, "sig")
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
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}
			if !resOut.Send(ctx, cfg) {
				return
			}
		}
	}, nil
}
