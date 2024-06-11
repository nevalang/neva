package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

var errIndexOutOfBounds = errors.New("index out of bounds")

type listAt struct{}

func (listAt) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
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
			var data []runtime.Msg
			var idx int64

			select {
			case <-ctx.Done():
				return
			case msg := <-dataIn:
				data = msg.List()
			}

			select {
			case <-ctx.Done():
				return
			case msg := <-idxIn:
				idx = msg.Int()
			}

			if idx < 0 || idx >= int64(len(data)) {
				select {
				case <-ctx.Done():
					return
				case errOut <- runtime.NewMapMsg(map[string]runtime.Msg{
					"text": runtime.NewStrMsg(errIndexOutOfBounds.Error()),
				}):
					continue
				}
			}

			select {
			case <-ctx.Done():
				return
			case resOut <- data[idx]:
			}
		}
	}, nil
}
