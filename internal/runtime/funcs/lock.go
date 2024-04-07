package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type lock struct{}

func (l lock) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, err := io.In.Port("sig")
	if err != nil {
		return nil, err
	}

	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Port("data")
	if err != nil {
		return nil, err
	}

	return l.Handle(sigIn, dataIn, dataOut), nil
}

func (lock) Handle(
	sigIn,
	dataIn,
	dataOut chan runtime.Msg,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-sigIn:
				select {
				case <-ctx.Done():
					return
				case dataMsg := <-dataIn:
					select {
					case <-ctx.Done():
						return
					case dataOut <- dataMsg:
					}
				}
			}
		}
	}
}
