package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intIsGreater struct{}

func (p intIsGreater) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
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
			actualMsg, ok := actualIn.Receive(ctx)
			if !ok {
				return
			}

			comparedMsg, ok := comparedIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(ctx, runtime.NewBoolMsg(actualMsg.Int() > comparedMsg.Int())) {
				return
			}
		}
	}, nil
}
