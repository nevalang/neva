package funcs

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

type osEnviron struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (o osEnviron) Create(io runtime.IO, _ messages.Msg) (func(ctx context.Context), error) {
	sigIn, err := io.In.Single("sig")
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
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			values := os.Environ()
			result := make([]messages.Msg, 0, len(values))
			for _, value := range values {
				result = append(result, messages.NewStringMsg(value))
			}

			if !resOut.Send(ctx, messages.NewListMsg(result)) {
				return
			}
		}
	}, nil
}
