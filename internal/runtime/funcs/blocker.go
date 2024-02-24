package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type blocker struct{}

func (l blocker) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	sigIn, err := io.In.Port("sig")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Port("data")
	if err != nil {
		return nil, err
	}

	return l.Handle(dataIn, sigIn, dataOut), nil
}

func (blocker) Handle(dataIn, sigIn, dataOut chan runtime.Msg) func(ctx context.Context) {
	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-sigIn:
				select {
				case <-ctx.Done():
					return
				case v := <-dataIn:
					select {
					case <-ctx.Done():
						return
					case dataOut <- v:
					}
				}
			}
		}
	}
}
