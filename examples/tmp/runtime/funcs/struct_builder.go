package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type structBuilder struct{}

func (s structBuilder) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	if len(io.In) == 0 {
		return nil, errors.New("cannot create struct builder without inports")
	}

	inports := make(map[string]chan runtime.Msg, len(io.In))
	for k, slots := range io.In {
		if len(slots) != 1 {
			return nil, errors.New("non-single port found: " + k)
		}
		inports[k] = slots[0]
	}

	vout, err := io.Out.Port("v")
	if err != nil {
		return nil, err
	}

	return s.Handle(inports, vout), nil
}

func (structBuilder) Handle(
	inports map[string]chan runtime.Msg,
	vout chan runtime.Msg,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				structure := make(map[string]runtime.Msg, len(inports))
				for name, inport := range inports {
					select {
					case <-ctx.Done():
						return
					case v := <-inport:
						structure[name] = v
					}
				}
				select {
				case <-ctx.Done():
					return
				case vout <- runtime.NewMapMsg(structure):
				}
			}
		}
	}
}
