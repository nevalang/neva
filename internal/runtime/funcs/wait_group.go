package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type waitGroup struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (g waitGroup) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	countIn, err := singleIn(io, "count")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	sigIn, err := singleIn(io, "sig")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	sigOut, err := singleOut(io, "sig")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
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
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			n, ok := countIn.Receive(ctx)
			if !ok {
				return
			}

			//nolint:intrange // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			for i := int64(0); i < n.Int(); i++ {
				if _, ok := sigIn.Receive(ctx); !ok {
					return
				}
			}

			if !sigOut.Send(ctx, emptyStruct()) {
				return
			}
		}
	}
}
