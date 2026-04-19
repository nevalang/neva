package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intMod struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (intMod) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	leftIn, err := io.In.Single("left") // numerator
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	denIn, err := io.In.Single("right") // denominator
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := io.Out.Single("res") // modulo
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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
