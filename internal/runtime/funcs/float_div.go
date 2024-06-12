package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

var errFloatDivideByZero = errors.New("float divide by zero")

type floatDiv struct{}

func (floatDiv) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
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
			var x, y float64
			select {
			case <-ctx.Done():
				return
			case msg := <-xIn:
				x = msg.Float()
			}
			select {
			case <-ctx.Done():
				return
			case msg := <-yIn:
				y = msg.Float()
			}
			if y == 0 {
				select {
				case <-ctx.Done():
					return
				case errOut <- runtime.NewMapMsg(map[string]runtime.Msg{
					"text": runtime.NewStrMsg(errFloatDivideByZero.Error()),
				}):
					continue
				}
			}
			select {
			case <-ctx.Done():
				return
			case resOut <- runtime.NewFloatMsg(x / y):
			}
		}
	}, nil
}
