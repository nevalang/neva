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
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-dataIn:
				var splitLst []runtime.Msg
				str := data.String()
				for i := 0; i < len(str); i++ {
					char := string(str[i])
					splitLst = append(splitLst, runtime.NewStrMsg(char))
				}
				resOut <- runtime.NewListMsg(splitLst...)
			}
		}
	}, nil
}
