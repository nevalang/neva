package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type streamToList struct{}

func (s streamToList) Create(
	io runtime.FuncIO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	seqIn, err := io.In.Port("seq")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		acc := []runtime.Msg{}

		for {
			var item map[string]runtime.Msg
			select {
			case <-ctx.Done():
				return
			case seqMsg := <-seqIn:
				item = seqMsg.Map()
			}

			acc = append(acc, item["data"])

			if item["last"].Bool() {
				select {
				case <-ctx.Done():
					return
				case resOut <- runtime.NewListMsg(acc...):
					acc = []runtime.Msg{} // reset
					continue
				}
			}
		}
	}, nil
}
