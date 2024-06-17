package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type listAt struct{}

func (listAt) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	idxIn, err := io.In.Single("idx")
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
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			idxMsg, ok := idxIn.Receive(ctx)
			if !ok {
				return
			}

			idx := idxMsg.Int()
			data := dataMsg.List()

			l := int64(len(data))
			if idx < -l || idx >= l {
				if !errOut.Send(ctx, errFromString("index out of bounds")) {
					return
				}
			}

			if idx < 0 {
				// support negative indexing:
				//	$l = [1, 2, 3]
				//	$l[-1] // 3
				idx += int64(len(data))
			}

			if !resOut.Send(ctx, data[idx]) {
				return
			}
		}
	}, nil
}
