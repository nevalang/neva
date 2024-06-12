package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intDiv struct{}

func (intDiv) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	xIn, err := io.In.Port("x")
	if err != nil {
		return nil, err
	}

	yIn, err := io.In.Port("y")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	errOut, err := io.Out.Port("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var x, y int64
			select {
			case <-ctx.Done():
				return
			case msg := <-xIn:
				x = msg.Int()
			}
			select {
			case <-ctx.Done():
				return
			case msg := <-yIn:
				y = msg.Int()
			}
			if y == 0 {
				select {
				case <-ctx.Done():
					return
				case errOut <- runtime.NewMapMsg(map[string]runtime.Msg{
					"text": runtime.NewStrMsg(errIntegerDivideByZero.Error()),
				}):
					continue
				}
			}
			select {
			case <-ctx.Done():
				return
			case resOut <- runtime.NewIntMsg(x / y):
			}
		}
	}, nil
}
