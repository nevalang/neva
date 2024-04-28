package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type list_iter struct{}

func (c list_iter) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	seqOut, err := io.Out.Port("seq")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case data, ok := <-dataIn:
				if !ok {
					return // lstIn channel closed
				}
				for i := 0; i < len(data.List()); i++ {
					select {
					case <-ctx.Done():
						return
					case seqOut <- data.List()[i]:
					}
				}
				select {
				case <-ctx.Done():
					return
				case seqOut <- nil:
				}
			case <-ctx.Done():
				return
			}
		}
	}, nil
}
