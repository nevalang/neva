package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intGreaterChecker struct{}

func (intGreaterChecker) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	aIn, err := io.In.Port("a")
	if err != nil {
		return nil, err
	}

	bIn, err := io.In.Port("b")
	if err != nil {
		return nil, err
	}

	yesOut, err := io.Out.Port("yes")
	if err != nil {
		return nil, err
	}

	noOut, err := io.Out.Port("no")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case a := <-aIn:
				select {
				case <-ctx.Done():
					return
				case b := <-bIn:
					if a.Int() > b.Int() {
						select {
						case <-ctx.Done():
							return
						case yesOut <- nil:
						}
						continue
					}
					select {
					case <-ctx.Done():
						return
					case noOut <- nil:
					}
				}
			}
		}
	}, nil
}
