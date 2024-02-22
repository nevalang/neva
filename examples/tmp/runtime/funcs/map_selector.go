package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type mapSelector struct{}

func (m mapSelector) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	// in
	vin, err := io.In.Port("v")
	if err != nil {
		return nil, err
	}
	kin, err := io.In.Port("k")
	if err != nil {
		return nil, err
	}

	// out
	okOut, err := io.Out.Port("ok")
	if err != nil {
		return nil, err
	}
	missOut, err := io.Out.Port("miss")
	if err != nil {
		return nil, err
	}

	// logic
	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case m := <-vin: // read map
				select {
				case <-ctx.Done():
					return
				case k := <-kin: // then read key
					var ( // figure out what and where to send
						msg runtime.Msg
						out chan runtime.Msg
					)
					v, ok := m.Map()[k.Str()]
					if ok {
						msg = v
						out = okOut
					} else {
						msg = k // if value not found, send missing key as a signal for miss outport
						out = missOut
					}
					select { // and send
					case <-ctx.Done():
						return
					case out <- msg:
						return
					}
				}
			}
		}
	}, nil
}
