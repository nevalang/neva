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
			acc     int64 = 0
			started bool  = false
		)

		for {
			var item map[string]runtime.Msg
			select {
			case <-ctx.Done():
				return
			case msg := <-seqIn:
				item = msg.Map()
			}

			if !started {
				acc = item["data"].Int()
				started = true
			} else {
				acc -= item["data"].Int()
			}

			if item["last"].Bool() {
				select {
				case <-ctx.Done():
					return
				case resOut <- runtime.NewIntMsg(acc):
					acc = 0
					started = false
					continue
				}
			}
		}
	}, nil
}
