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

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			if !resOut.Send(ctx, argsListMsg(os.Args)) {
				return
			}
		}
	}, nil
}

// argsListMsg converts argv list to runtime list message.
func argsListMsg(argv []string) runtime.ListMsg {
	result := make([]runtime.Msg, len(argv))
	for i := range argv {
		result[i] = runtime.NewStringMsg(argv[i])
	}
	return runtime.NewListMsg(result)
}
