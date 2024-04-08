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
		var listMsg, idxMsg runtime.Msg

		for {
			select {
			case <-ctx.Done():
				return
			case listMsg = <-listIn:
			}

			select {
			case <-ctx.Done():
				return
			case idxMsg = <-indexIn:
			}

			idx := idxMsg.Int()
			list := listMsg.List()

			if idx < 0 || idx >= int64(len(list)) {
				select {
				case <-ctx.Done():
					return
				case errOut <- runtime.NewStrMsg("Index out of bounds"):
					continue
				}
			}

			select {
			case <-ctx.Done():
				return
			case resOut <- list[idx]:
			}
		}
	}, nil
}
