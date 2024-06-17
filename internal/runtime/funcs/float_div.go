package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

var errFloatDivideByZero = errors.New("float divide by zero")

type floatDiv struct{}

func (floatDiv) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	xIn, err := io.In.Single("x")
	if err != nil {
		return nil, err
	}

	yIn, err := io.In.Single("y")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.SingleOutport("res")
	if err != nil {
		return nil, err
	}

	errOut, err := io.Out.SingleOutport("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			xMsg, ok := xIn.Receive(ctx)
			if !ok {
				return
			}

			yMsg, ok := yIn.Receive(ctx)
			if !ok {
				return
			}

			if yMsg.Float() == 0 {
				if !errOut.Send(ctx, errFromString("divide by zero")) {
					return
				}
				continue
			}

			if !resOut.Send(
				ctx,
				runtime.NewFloatMsg(
					xMsg.Float()/yMsg.Float(),
				),
			) {
				return
			}
		}
	}, nil
}
