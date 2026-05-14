package funcs

import (
	"context"
	"fmt"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type panicker struct{}

func (p panicker) Create(
	runtimeIO runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	msgIn, err := runtimeIO.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		panicMsg, ok := msgIn.Receive(ctx)
		if !ok {
			return
		}

		if _, err := fmt.Fprintln(os.Stderr, "panic:", panicMsg); err != nil {
			panic(err)
		}

		writeTerminationTrace("panic cause dataflow trace", runtimeIO, panicMsg)

		runtime.Terminate(ctx, 1)
	}, nil
}
