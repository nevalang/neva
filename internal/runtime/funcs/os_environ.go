package funcs

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type osEnviron struct{}

func (o osEnviron) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := io.In.Single("sig")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			values := os.Environ()
			result := make([]runtime.Msg, 0, len(values))
			for _, value := range values {
				result = append(result, runtime.NewStringMsg(value))
			}

			if !resOut.Send(ctx, runtime.NewListMsg(result)) {
				return
			}
		}
	}, nil
}
