package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type fanIn struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (fanIn) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	data, err := arrayIn(io, "data")
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
			dataMsg, ok := data.Select(ctx)
			if !ok {
				return
			}

			if !resOut.Send(ctx, dataMsg) {
				return
			}
		}
	}, nil
}
