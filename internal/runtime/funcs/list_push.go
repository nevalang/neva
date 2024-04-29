package funcs

import (
	"context"
	"slices"

	"github.com/nevalang/neva/internal/runtime"
)

type listPush struct{}

func (p listPush) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
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
		var (
			data runtime.Msg
			lst  runtime.Msg
		)

		for {
			select {
			case <-ctx.Done():
				return
			case data = <-dataIn:
			}

			select {
			case <-ctx.Done():
				return
			case lst = <-lstIn:
			}

			lstCopy := slices.Clone(lst.List())
			lstCopy = append(lstCopy, data)

			select {
			case <-ctx.Done():
				return
			case resOut <- runtime.NewListMsg(lstCopy...):
			}
		}
	}, nil
}
