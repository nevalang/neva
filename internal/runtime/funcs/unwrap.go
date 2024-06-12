package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type unwrap struct{}

func (unwrap) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.SingleInport("data")
	if err != nil {
		return nil, err
	}

	someOut, err := io.Out.SingleOutport("some")
	if err != nil {
		return nil, err
	}

	noneOut, err := io.Out.SingleOutport("none")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var dataMsg runtime.Msg

		for {
			select {
			case <-ctx.Done():
				return
			case dataMsg = <-dataIn:
			}

			if dataMsg == nil {
				select {
				case <-ctx.Done():
					return
				case noneOut <- runtime.NewMapMsg(nil):
				}
				continue
			}

			select {
			case <-ctx.Done():
				return
			case someOut <- dataMsg:
			}
		}
	}, nil
}
