package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type floatIsGreater struct{}

func (p floatIsGreater) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	actualIn, err := io.In.Single("left")
	if err != nil {
		return nil, err
	}

	comparedIn, err := io.In.Single("right")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
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
				runtime.NewBoolMsg(val1.Float() > val2.Float()),
			) {
				return
			}
		}
	}, nil
}
