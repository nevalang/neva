package funcs

import (
	"context"
	"strconv"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type formatBool struct{}

func (formatBool) Create(
	//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	io runtime.IO,
	_ messages.Msg,
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

			res := strconv.FormatBool(data.Bool())
			if !resOut.Send(ctx, messages.NewStringMsg(res)) {
				return
			}
		}
	}, nil
}
