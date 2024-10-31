package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intMod struct{}

func (intMod) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	numIn, err := io.In.Single("acc") // numerator
	if err != nil {
		return nil, err
	}

	denIn, err := io.In.Single("el") // denominator
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			numMsg, ok := numIn.Receive(ctx)
			if !ok {
				return
			}

			denMsg, ok := denIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(
				ctx,
				runtime.NewIntMsg(
					numMsg.Int()%denMsg.Int(),
				),
			) {
				return
			}
		}
	}, nil
}
