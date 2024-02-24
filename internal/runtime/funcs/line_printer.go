package funcs

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type linePrinter struct{}

func (p linePrinter) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}
	dataOut, err := io.Out.Port("sig")
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case data := <-dataIn:
				select {
				case <-ctx.Done():
					return
				default:
					fmt.Println(data)
					select {
					case <-ctx.Done():
						return
					case dataOut <- data:
					}
				}
			}
		}
	}, nil
}
