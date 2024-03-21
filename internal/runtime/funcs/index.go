package funcs

import (
	"context"
	"github.com/nevalang/neva/internal/runtime"
)

type index struct{}

func (p index) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	listIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	indexIn, err := io.In.Port("idx")
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
			select {
			case <-ctx.Done():
				return
			case data := <-listIn:
				select {
				case <-ctx.Done():
				case idx := <-indexIn:
					select {
					case <-ctx.Done():
					default:
						intIdx := idx.Int()
						lst := data.List()
						if intIdx < 0 || intIdx >= int64(len(lst)) {
							errOut <- runtime.NewStrMsg("Index out of bounds")
						} else {
							resOut <- lst[intIdx]
						}
					}
				}
			}
		}
	}, nil
}
