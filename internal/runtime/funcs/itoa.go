package funcs

import (
	"context"
	"strconv"

	"github.com/nevalang/neva/internal/runtime"
)

type itoa struct{}

func (itoa) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := io.Out.Single("res")
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

			res := strconv.FormatInt(data.Int(), 10)
			if !resOut.Send(ctx, runtime.NewStringMsg(res)) {
				return
			}
		}
	}, nil
}
