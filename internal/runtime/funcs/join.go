package funcs

import (
	"context"
	"github.com/nevalang/neva/internal/runtime"
	"strings"
)

type join struct{}

func (p join) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}
	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-dataIn:
				var str strings.Builder
				for i := 0; i < len(data.List()); i++ {
					str.WriteString(data.List()[i].String())
				}
				resOut <- runtime.NewStrMsg(str.String())
			}
		}
	}, nil
}
