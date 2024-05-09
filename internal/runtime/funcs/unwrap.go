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
		var dataMsg runtime.Msg

		for {
			var item map[string]runtime.Msg
			select {
			case <-ctx.Done():
				return
			case dataMsg = <-dataIn:
				item = dataMsg.Map()
			}
			if item["last"].Bool() {
				select {
				case someOut <- item["data"]:
					noneOut <- runtime.NewMapMsg(nil)
					return
				case <-ctx.Done():
					return
				}
			} else {
				select {
				case <-ctx.Done():
					return
				case someOut <- item["data"]:
				}
			}
		}
	}, nil
}
