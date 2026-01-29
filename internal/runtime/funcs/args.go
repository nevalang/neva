package funcs

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type args struct{}

func (a args) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := io.In.Single("sig")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Single("data")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		if _, ok := sigIn.Receive(ctx); !ok {
			return
		}

		result := make([]runtime.Msg, len(os.Args))
		for i := range os.Args {
			result = append(result, runtime.NewStringMsg(os.Args[i]))
		}

		if !dataOut.Send(ctx, runtime.NewListMsg(result)) {
			return
		}
	}, nil
}
