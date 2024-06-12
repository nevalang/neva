package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"golang.org/x/exp/slices"
)

type listSortString struct{}

func (p listSortString) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.SingleInport("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.SingleOutport("res")
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

			clone := slices.Clone(data.List())
			slices.SortFunc(clone, func(i, j runtime.Msg) bool {
				return i.Str() < j.Str()
			})

			select {
			case <-ctx.Done():
				return
			case resOut <- runtime.NewListMsg(clone...):
			}
		}
	}, nil
}
