package funcs

import (
	"context"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

type timeAfter struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (timeAfter) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	durIn, err := singleIn(io, "dur")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	sigOut, err := singleOut(io, "sig")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			durMsg, ok := durIn.Receive(ctx)
			if !ok {
				return
			}

			time.Sleep(time.Duration(durMsg.Int()))

			if !sigOut.Send(ctx, emptyStruct()) {
				return
			}
		}
	}, nil
}
