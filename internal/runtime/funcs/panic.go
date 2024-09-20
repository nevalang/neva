package funcs

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type panicker struct{}

func (p panicker) Create(
	io runtime.IO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	msgIn, err := io.In.Single("msg")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		panicMsg, ok := msgIn.Receive(ctx)
		if !ok {
			return
		}

		cancel := ctx.Value("cancel").(context.CancelFunc)
		cancel()

		fmt.Printf("panic: %v\n", panicMsg) // stderr?
	}, nil
}
