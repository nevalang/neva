package funcs

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type println struct{}

func (p println) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	sigOut, err := io.Out.Port("sig")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var data runtime.Msg

		for {
			select {
			case <-ctx.Done():
				return
			case data = <-dataIn:
			}

			if _, err := fmt.Println(data); err != nil {
				panic(err)
			}

			select {
			case <-ctx.Done():
				return
			case sigOut <- data:
			}
		}
	}, nil
}
