package funcs

import (
	"context"
	"unicode/utf8"

	"github.com/nevalang/neva/internal/runtime"
)

type stringAt struct{}

func (stringAt) Create(io runtime.IO, _ runtime.Msg) (func(context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	idxIn, err := io.In.Single("idx")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	errOut, err := io.Out.Single("err")
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
			data := dataMsg.Str()
			l := int64(utf8.RuneCountInString(data))

			if idx < -l || idx >= l {
				if !errOut.Send(ctx, errFromString("index out of bounds")) {
					return
				}
			}

			for i, r := range data {
				if int64(i) == idx {
					if !resOut.Send(ctx, runtime.NewStrMsg(string(r))) {
						return
					}
					break
				}
			}
		}
	}, nil
}
