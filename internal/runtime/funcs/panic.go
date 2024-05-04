package funcs

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type panicker struct{}

func (p panicker) Create(
	io runtime.FuncIO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	msgIn, err := io.In.Port("msg")
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		case panicMsg := <-msgIn:
			cancel := ctx.Value("cancel").(context.CancelFunc)
			cancel()
			fmt.Printf("panic: %v\n", panicMsg)
		}
	}, nil
}
