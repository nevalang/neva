//nolint:dupl
package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intIsLesser struct{}

//nolint:varnamelen
func (p intIsLesser) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	actualIn, err := io.In.Single("left")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	comparedIn, err := io.In.Single("right")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			//nolint:varnamelen
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
