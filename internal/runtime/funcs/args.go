package funcs

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type args struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (a args) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
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
		if _, ok := sigIn.Receive(ctx); !ok {
			return
		}

		result := make([]runtime.Msg, len(os.Args))
		for i := range os.Args {
			//nolint:makezero // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			result = append(result, runtime.NewStringMsg(os.Args[i]))
		}

		if !resOut.Send(ctx, runtime.NewListMsg(result)) {
			return
		}
	}, nil
}
