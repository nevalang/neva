package funcs

import (
	"context"
	"github.com/nevalang/neva/internal/runtime"
)

type and struct{}

func (p and) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	AIn, err := io.In.Port("A")
	if err != nil {
		return nil, err
	}
	BIn, err := io.In.Port("B")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var (
			AVAL runtime.Msg
			BVAL runtime.Msg
		)

		for {
			select {
			case <-ctx.Done():
				return
			case AVAL = <-AIn:
			}

			select {
			case <-ctx.Done():
				return
			case BVAL = <-BIn:
			}

			select {
			case <-ctx.Done():
				return

			default:
				resOut <- runtime.NewBoolMsg(BVAL.Bool() && AVAL.Bool())
			}
		}
	}, nil
}
