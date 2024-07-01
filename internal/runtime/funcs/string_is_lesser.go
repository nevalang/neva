package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type strIsLesser struct{}

func (p strIsLesser) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	actualIn, err := io.In.Single("actual")
	if err != nil {
		return nil, err
	}

	comparedIn, err := io.In.Single("compared")
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

			if !resOut.Send(ctx, runtime.NewBoolMsg(val1.Str() < val2.Str())) {
				return
			}
		}
	}, nil
}
