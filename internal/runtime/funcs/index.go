package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type index struct{}

func (p index) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	listIn, err := io.In.SingleInport("data")
	if err != nil {
		return nil, err
	}

	indexIn, err := io.In.SingleInport("idx")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.SingleOutport("res")
	if err != nil {
		return nil, err
	}

	errOut, err := io.Out.SingleOutport("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			listMsg, ok := listIn.Receive(ctx)
			if !ok {
				return
			}

			idxMsg, ok := indexIn.Receive(ctx)
			if !ok {
				return
			}

			idx := idxMsg.Int()
			list := listMsg.List()

			if idx < 0 || idx >= int64(len(list)) {
				if !errOut.Send(ctx, errFromString("Index out of bounds")) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, list[idx]) {
				return
			}
		}
	}, nil
}
