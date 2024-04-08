package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type split struct{}

func (p split) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}
	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var data runtime.Msg

		for {
			select {
			case <-ctx.Done():
				return
			case data = <-dataIn:
			}

			str := data.Str()
			splitLst := make([]runtime.Msg, 0, len(str))
			for _, r := range str {
				splitLst = append(splitLst, runtime.NewStrMsg(string(r)))
			}

			select {
			case <-ctx.Done():
				return
			case resOut <- runtime.NewListMsg(splitLst...):
			}
		}
	}, nil
}
