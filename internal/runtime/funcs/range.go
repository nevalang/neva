package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type ranger struct{}

func (ranger) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	fromIn, err := io.In.Port("from")
	if err != nil {
		return nil, err
	}

	toIn, err := io.In.Port("to")
	if err != nil {
		return nil, err
	}

	dataOut, err := io.Out.Port("data")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var fromMsg, toMsg runtime.Msg

		for {
			select {
			case <-ctx.Done():
				return
			case fromMsg = <-fromIn:
			}

			select {
			case <-ctx.Done():
				return
			case toMsg = <-toIn:
			}

			for i := fromMsg.Int(); i < toMsg.Int(); i++ {
				select {
				case <-ctx.Done():
					return
				case dataOut <- runtime.NewIntMsg(i):
				}
			}

			select {
			case <-ctx.Done():
				return
			case dataOut <- nil:
			}
		}
	}, nil
}
