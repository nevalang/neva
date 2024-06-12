package funcs

import (
	"context"

	"golang.org/x/exp/utf8string"

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
			var data *utf8string.String
			var idx int64

			select {
			case <-ctx.Done():
				return
			case msg := <-dataIn:
				data = utf8string.NewString(msg.Str())
			}

			select {
			case <-ctx.Done():
				return
			case msg := <-idxIn:
				idx = msg.Int()
			}

			if idx < 0 || idx >= int64(data.RuneCount()) {
				select {
				case <-ctx.Done():
					return
				case errOut <- errorFromString(errIndexOutOfBounds.Error()):
					continue
				}
			}

			select {
			case <-ctx.Done():
				return
			case resOut <- runtime.NewStrMsg(string(data.At(int(idx)))):
			}
		}
	}, nil
}
