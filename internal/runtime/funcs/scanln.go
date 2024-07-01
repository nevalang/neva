package funcs

import (
	"bufio"
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type scanln struct{}

func (r scanln) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := io.In.Single("sig")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Single("data")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			reader := bufio.NewReader(os.Stdin)

			if _, ok := sigIn.Receive(ctx); !ok {
				return
			}

			bb, _, err := reader.ReadLine()
			if err != nil {
				panic(err)
			}

			if !dataOut.Send(ctx, runtime.NewStrMsg(string(bb))) {
				return
			}
		}
	}, nil
}
