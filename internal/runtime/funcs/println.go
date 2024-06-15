package funcs

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type println struct{}

func (p println) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	sigOut, err := io.Out.SingleOutport("sig")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			if _, err := fmt.Println(data); err != nil {
				panic(err)
			}

			if !sigOut.Send(ctx, data) {
				return
			}
		}
	}, nil
}
