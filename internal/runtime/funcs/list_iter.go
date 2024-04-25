package funcs

import (
	"context"
	"github.com/nevalang/neva/internal/runtime"
)

type list_iter struct{}

func (c list_iter) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	lstIn, err := io.In.Port("lst")
	if err != nil {
		return nil, err
	}
	outport, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case data, ok := <-lstIn:
				if !ok {
					return // lstIn channel closed
				}
				for i := 0; i < len(data.List()); i++ {
					select {
					case <-ctx.Done():
						return
					case outport <- data.List()[i]:
					}
				}
				select {
				case <-ctx.Done():
					return
				case outport <- nil:
				}
			case <-ctx.Done():
				return
			}
		}
	}, nil
}
