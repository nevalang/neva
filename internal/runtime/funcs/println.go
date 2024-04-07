package funcs

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type println struct{}

func (p println) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	sigOut, err := io.Out.Port("sig")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case dataMsg := <-dataIn:
				select {
				case <-ctx.Done():
					return
				default:
					_, err := fmt.Println(dataMsg)
					if err != nil {
						panic(err)
					}
					select {
					case <-ctx.Done():
						return
					case sigOut <- dataMsg:
					}
				}
			}
		}
	}, nil
}
