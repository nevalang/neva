package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type blocker struct{}

func (l blocker) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	vin, err := io.In.Port("v")
	if err != nil {
		return nil, err
	}

	sig, err := io.In.Port("sig")
	if err != nil {
		return nil, err
	}

	vout, err := io.Out.Port("v")
	if err != nil {
		return nil, err
	}

	return l.Handle(vin, sig, vout), nil
}

func (blocker) Handle(vin, sig, vout chan runtime.Msg) func(ctx context.Context) {
	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-sig:
				select {
				case <-ctx.Done():
					return
				case v := <-vin:
					select {
					case <-ctx.Done():
						return
					case vout <- v:
					}
				}
			}
		}
	}
}
