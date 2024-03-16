package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type intEq struct{}

func (intEq) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	aIn, err := io.In.Port("a")
	if err != nil {
		return nil, err
	}

	bIn, err := io.In.Port("b")
	if err != nil {
		return nil, err
	}

	thenOut, err := io.Out.Port("then")
	if err != nil {
		return nil, err
	}

	elseOut, err := io.Out.Port("else")
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
					if a.Int() == b.Int() {
						select {
						case <-ctx.Done():
							return
						case thenOut <- a:
						}
						continue
					}
					elseMsg := runtime.NewMapMsg(map[string]runtime.Msg{
						"a": a,
						"b": b,
					})
					select {
					case <-ctx.Done():
						return
					case elseOut <- elseMsg:
					}
				}
			}
		}
	}, nil
}
