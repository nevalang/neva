package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type stringAt struct{}

func (stringAt) Create(io runtime.FuncIO, _ runtime.Msg) (func(context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	idxIn, err := io.In.Port("idx")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	errOut, err := io.Out.Port("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var data string
			var idx int64

			select {
			case <-ctx.Done():
				return
			case msg := <-dataIn:
				data = msg.Str()
			}

			select {
			case <-ctx.Done():
				return
			case msg := <-idxIn:
				idx = msg.Int()
			}

			if idx >= 0 && idx < int64(len(data)) {
				var res rune
				var found bool
				for i, r := range data {
					if int64(i) == idx {
						res = r
						found = true
						break
					}
				}
				if found {
					select {
					case <-ctx.Done():
						return
					case resOut <- runtime.NewStrMsg(string(data.At(int(idx)))):
						continue
					}
				}
			}

			select {
			case <-ctx.Done():
				return
			case errOut <- errorFromString(errIndexOutOfBounds.Error()):
			}
		}
	}, nil
}
