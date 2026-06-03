package funcs

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type stringsToLower struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (p stringsToLower) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
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

			result := strings.ToLower(data.Str())
			if !resOut.Send(ctx, runtime.NewStringMsg(result)) {
				return
			}
		}
	}, nil
}
