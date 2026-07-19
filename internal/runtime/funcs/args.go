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
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			if !resOut.Send(ctx, runtime.NewListStringMsg(os.Args)) {
				return
			}
		}
	}, nil
}
