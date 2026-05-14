package funcs

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type osExit struct{}

func (osExit) Create(runtimeIO runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	codeIn, err := runtimeIO.In.Single("code")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		codeMsg, ok := codeIn.Receive(ctx)
		if !ok {
			return
		}

		exitCode := codeMsg.Int()
		if exitCode < 0 {
			panic(fmt.Sprintf("runtime invariant: os.exit code must be non-negative, got %d", exitCode))
		}

		runtime.Terminate(ctx, int(exitCode))
	}, nil
}
