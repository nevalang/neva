package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type ternary struct{}

func (p ternary) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	ifIn, err := io.In.Single("if")
	if err != nil {
		return nil, err
	}

	thenIn, err := io.In.Single("then")
	if err != nil {
		return nil, err
	}

	elseIn, err := io.In.Single("else")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
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
