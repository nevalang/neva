package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamJust struct{}

func (streamJust) Create(input runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := input.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := input.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(ctx, streamItem(data, 0, true)) {
				return
			}
		}
	}, nil
}
