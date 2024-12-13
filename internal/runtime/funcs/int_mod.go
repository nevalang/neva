package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intMod struct{}

func (intMod) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	leftIn, err := io.In.Single("left") // numerator
	if err != nil {
		return nil, err
	}

	denIn, err := io.In.Single("right") // denominator
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res") // modulo
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			numMsg, ok := leftIn.Receive(ctx)
			if !ok {
				return
			}

			denMsg, ok := denIn.Receive(ctx)
			if !ok {
				return
			}

			num := numMsg.Int()
			den := denMsg.Int()
			if !resOut.Send(ctx, runtime.NewIntMsg(num%den)) {
				return
			}
		}
	}, nil
}
