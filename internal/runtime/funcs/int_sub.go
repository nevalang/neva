package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intSub struct{}

func (intSub) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	seqIn, err := io.In.Port("seq")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var (
			acc        int64 = 0
			inProgress bool  = false
			cur        runtime.Msg
		)

		for {
			select {
			case <-ctx.Done():
				return
			case cur = <-seqIn:
			}

			if cur == nil {
				select {
				case <-ctx.Done():
					return
				case resOut <- runtime.NewIntMsg(acc):
					acc = 0
					inProgress = false
					continue
				}
			}

			if !inProgress {
				acc = cur.Int()
				inProgress = true
			} else {
				acc -= cur.Int()
			}
		}
	}, nil
}
