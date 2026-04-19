package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type ternarySelector struct{}

//nolint:gocognit,varnamelen
func (p ternarySelector) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	ifIn, err := io.In.Single("if")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	thenIn, err := io.In.Single("then")
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	elseIn, err := io.In.Single("else")
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
			dataMsg, ok := ifIn.Receive(ctx)
			if !ok {
				return
			}

			thenMsg, ok := thenIn.Receive(ctx)
			if !ok {
				return
			}

			elseMsg, ok := elseIn.Receive(ctx)
			if !ok {
				return
			}

			var resMsg runtime.Msg
			if dataMsg.Bool() {
				resMsg = thenMsg
			} else {
				resMsg = elseMsg
			}

			if !resOut.Send(ctx, resMsg) {
				return
			}
		}
	}, nil
}
