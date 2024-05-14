package funcs

import (
	"context"
	"github.com/nevalang/neva/internal/runtime"
)

type intIsGreater struct{}

func (p intIsGreater) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	actualIn, err := io.In.Port("actual")
	if err != nil {
		return nil, err
	}
	comparedIn, err := io.In.Port("compared")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var (
			val1 runtime.Msg
			val2 runtime.Msg
		)

		for {
			select {
			case <-ctx.Done():
				return
			case val1 = <-actualIn:
			}

			select {
			case <-ctx.Done():
				return
			case val2 = <-comparedIn:
			}

			select {
			case <-ctx.Done():
				return
			case resOut <- runtime.NewBoolMsg(val1.Int() > val2.Int()):
			}
		}
	}, nil
}
