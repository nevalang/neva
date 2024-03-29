package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intMultiplier struct{}

func (intMultiplier) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	streamIn, err := io.In.Port("stream")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var res int64 = 1
		for {
			select {
			case <-ctx.Done():
				return
			case streamItem := <-streamIn:
				if streamItem == nil {
					select {
					case <-ctx.Done():
						return
					case resOut <- runtime.NewIntMsg(res):
						continue
					}
				}
				res *= streamItem.Int()
			}
		}
	}, nil
}
