package funcs

import (
	"context"
	"github.com/nevalang/neva/internal/runtime"
)

type intSubtracter struct{}

func (intSubtracter) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	streamIn, err := io.In.Port("stream")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		flag := false
		var res int64
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
						res = 0
						continue
					}
				}
				if !flag {
					res = streamItem.Int()
					flag = true
				} else {
					res -= streamItem.Int()
				}
			}
		}
	}, nil
}
