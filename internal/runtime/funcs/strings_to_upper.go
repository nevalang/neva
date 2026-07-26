package funcs

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type stringsToUpper struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (p stringsToUpper) Create(io runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
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

			result := strings.ToUpper(data.Str())
			if !resOut.Send(ctx, messages.NewStringMsg(result)) {
				return
			}
		}
	}, nil
}
