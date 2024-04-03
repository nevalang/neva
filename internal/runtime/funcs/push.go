package funcs

import (
	"context"
	"github.com/nevalang/neva/internal/runtime"
)

type listpush struct{}

func (p listpush) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}
	lstIn, err := io.In.Port("lst")
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
				select {
				case <-ctx.Done():
					return
				case lst := <-lstIn:
					newLst := lst.List()
					newLst = append(newLst, data)
					resOut <- runtime.NewListMsg(newLst...)
				}
			}
		}
	}, nil
}
