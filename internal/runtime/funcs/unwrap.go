package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type unwrap struct{}

func (unwrap) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	someOut, err := io.Out.Port("some")
	if err != nil {
		return nil, err
	}

	noneOut, err := io.Out.Port("none")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case dataMsg := <-dataIn:
				select {
				case <-ctx.Done():
					return
				default:
					if dataMsg != nil {
						select {
						case <-ctx.Done():
							return
						case someOut <- dataMsg:
						}
					} else {
						select {
						case <-ctx.Done():
							return
						case noneOut <- dataMsg:
						}
					}
				}
			}
		}
	}, nil
}
