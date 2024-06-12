package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type structBuilder struct{}

func (s structBuilder) Create(
	io runtime.FuncIO,
	_ runtime.Msg,
) (func(ctx context.Context), error) {
	if len(io.In) == 0 {
		return nil, errors.New("cannot create struct builder without inports")
	}

	inports := make(map[string]chan runtime.Msg, len(io.In))
	for inportName, inportSlots := range io.In {
		if len(inportSlots) != 1 {
			return nil, errors.New("non-single port found: " + inportName)
		}
		inports[inportName] = inportSlots[0]
	}

	outport, err := io.Out.SingleOutport("msg")
	if err != nil {
		return nil, err
	}

	return s.Handle(inports, outport), nil
}

func (structBuilder) Handle(
	inports map[string]chan runtime.Msg,
	outport chan runtime.Msg,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		for {
			var structure = make(map[string]runtime.Msg, len(inports))

			for inportName, inportChan := range inports {
				select {
				case <-ctx.Done():
					return
				case msg := <-inportChan:
					structure[inportName] = msg
				}
			}

			select {
			case <-ctx.Done():
				return
			case outport <- runtime.NewMapMsg(structure):
			}
		}
	}
}
