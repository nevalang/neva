package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type waitGroup struct{}

func (g waitGroup) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	countIn, err := io.In.Single("count")
	if err != nil {
		return nil, err
	}

	sigIn, err := io.In.Single("sig")
	if err != nil {
		return nil, err
	}

	sigOut, err := io.Out.Single("sig")
	if err != nil {
		return nil, err
	}

	return g.Handle(countIn, sigIn, sigOut), nil
}

func (waitGroup) Handle(
	countIn runtime.SingleInport,
	sigIn runtime.SingleInport,
	sigOut runtime.SingleOutport,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		for {
			n, ok := countIn.Receive(ctx)
			if !ok {
				return
			}

			for i := int64(0); i < n.Int(); i++ {
				if _, ok := sigIn.Receive(ctx); !ok {
					return
				}
			}

			if !sigOut.Send(ctx, nil) {
				return
			}
		}
	}
}
