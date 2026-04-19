package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type listToStream struct{}

func (c listToStream) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
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
			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			list := data.List()

			for idx := range list {
				item := streamItem(
					list[idx],
					int64(idx),
					idx == len(list)-1,
				)

				if !resOut.Send(ctx, item) {
					return
				}
			}
		}
	}, nil
}
