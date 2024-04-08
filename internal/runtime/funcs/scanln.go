package funcs

import (
	"bufio"
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type scanln struct{}

func (r scanln) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := io.In.Port("sig")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Port("data")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var reader = bufio.NewReader(os.Stdin)

		for {
			select {
			case <-ctx.Done():
				return
			case <-sigIn:
			}

			bb, _, err := reader.ReadLine()
			if err != nil {
				panic(err)
			}

			select {
			case <-ctx.Done():
				return
			case dataOut <- runtime.NewStrMsg(string(bb)):
			}
		}
	}, nil
}
