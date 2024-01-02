package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type addInts struct{}

func (addInts) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	ain, err := io.In.Port("a")
	if err != nil {
		return nil, err
	}

	bin, err := io.In.Port("b")
	if err != nil {
		return nil, err
	}

	vout, err := io.Out.Port("v")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case v1 := <-ain:
				select {
				case <-ctx.Done():
					return
				case v2 := <-bin:
					select {
					case <-ctx.Done():
						return
					case vout <- runtime.NewIntMsg(v1.Int() + v2.Int()):
					}
				}
			}
		}
	}, nil
}

type addFloats struct{}

func (a addFloats) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	ain, err := io.In.Port("a")
	if err != nil {
		return nil, err
	}
	bin, err := io.In.Port("b")
	if err != nil {
		return nil, err
	}
	vout, err := io.Out.Port("v")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case v1 := <-ain:
				select {
				case <-ctx.Done():
					return
				case v2 := <-bin:
					select {
					case <-ctx.Done():
						return
					case vout <- runtime.NewFloatMsg(v1.Float() + v2.Float()):
					}
				}
			}
		}
	}, nil
}
