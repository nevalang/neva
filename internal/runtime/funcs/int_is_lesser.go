//nolint:dupl // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intIsLesser struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (p intIsLesser) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	actualIn, err := singleIn(io, "left")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	comparedIn, err := singleIn(io, "right")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := singleOut(io, "res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
			actualMsg, ok := actualIn.Receive(ctx)
			if !ok {
				return
			}

			comparedMsg, ok := comparedIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(ctx, runtime.NewBoolMsg(actualMsg.Int() < comparedMsg.Int())) {
				return
			}
		}
	}, nil
}
