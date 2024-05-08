package funcs

import (
	"context"
	"github.com/nevalang/neva/internal/runtime"
)

type assert struct{}

func (p assert) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	valIn, err := io.In.Port("val")
	if err != nil {
		return nil, err
	}
	val2In, err := io.In.Port("val2")
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
			case val1 = <-valIn:
			}

			select {
			case <-ctx.Done():
				return
			case val2 = <-val2In:
			}

			select {
			case <-ctx.Done():
				return
			case resOut <- runtime.NewBoolMsg(val1 == val2):
			}
		}
	}, nil
}
