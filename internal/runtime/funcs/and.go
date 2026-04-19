//nolint:dupl
package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type and struct{}

//nolint:varnamelen
func (p and) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	aIn, err := io.In.Single("left")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	bIn, err := io.In.Single("right")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	// TODO send false as soon as A in is false, but do it correctly
	return func(ctx context.Context) {
		for {
			//nolint:varnamelen
			aMsg, ok := aIn.Receive(ctx)
			if !ok {
				return
			}

			bMsg, ok := bIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(
				ctx,
				runtime.NewBoolMsg(aMsg.Bool() && bMsg.Bool()),
			) {
				return
			}
		}
	}, nil
}
