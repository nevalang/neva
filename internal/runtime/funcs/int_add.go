package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intAdd struct{}

func (intAdd) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	seqIn, err := io.In.Port("seq")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var res int64
		for {
			select {
			case <-ctx.Done():
				return
			case streamItem := <-seqIn:
				if streamItem == nil {
					select {
					case <-ctx.Done():
						return
					case resOut <- runtime.NewIntMsg(res):
						res = 0
						continue
					}
				}
				res += streamItem.Int()
			}
		}
	}, nil
}
