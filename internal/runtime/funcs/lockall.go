package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type lockAll struct{}

func (l lockAll) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	sigIn, ok := io.In["sig"]
	if !ok {
		return nil, errors.New("inport 'sig' is required")
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

func (lockAll) Handle(
	sigIn []chan runtime.Msg,
	dataIn,
	dataOut chan runtime.Msg,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		var data runtime.Msg

		for {
			for _, sig := range sigIn {
				select {
				case <-ctx.Done():
					return
				case <-sig:
				}
			}

			select {
			case <-ctx.Done():
				return
			case data = <-dataIn:
			}

			select {
			case <-ctx.Done():
				return
			case dataOut <- data:
			}
		}
	}
}
