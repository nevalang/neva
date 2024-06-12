package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type eq struct{}

func (p eq) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	actualIn, err := io.In.SingleInport("actual")
	if err != nil {
		return nil, err
	}

	comparedIn, err := io.In.SingleInport("compared")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.SingleOutport("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			val1, ok := actualIn.Receive(ctx)
			if !ok {
				return
			}

			val2, ok := comparedIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(
				ctx,
				runtime.NewBoolMsg(val1 == val2),
			) {
				return
			}
		}
	}, nil
}
