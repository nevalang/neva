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
	if len(io.In.Ports()) == 0 {
		return nil, errors.New("cannot create struct builder without inports")
	}

	inports := make(map[string]runtime.SingleInport, len(io.In.Ports()))
	for inportName, inportSlots := range io.In.Ports() {
		if inportSlots.Single() == nil {
			return nil, errors.New("non-single port found: " + inportName)
		}
		inports[inportName] = *inportSlots.Single()
	}

	outport, err := io.Out.SingleOutport("msg")
	if err != nil {
		return nil, err
	}

	return s.Handle(inports, outport), nil
}

func (structBuilder) Handle(
	inports map[string]runtime.SingleInport,
	outport runtime.SingleOutport,
) func(ctx context.Context) {
	return func(ctx context.Context) {
		for {
			var structure = make(map[string]runtime.Msg, len(inports))

			for inportName, inportChan := range inports {
				msg, ok := inportChan.Receive(ctx)
				if !ok {
					return
				}
				structure[inportName] = msg
			}

			if !outport.Send(ctx, runtime.NewMapMsg(structure)) {
				return
			}
		}
	}
}
