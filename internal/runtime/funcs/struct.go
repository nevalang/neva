package funcs

import (
	"context"
	"errors"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

type structBuilder struct{}

func (s structBuilder) Create(
	io runtime.IO,
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

	outport, err := io.Out.Single("res")
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
			names := make([]string, 0, len(inports))
			fields := make([]runtime.Msg, 0, len(inports))
			var mu sync.Mutex
			var wg sync.WaitGroup
			for inportName, inportChan := range inports {
				wg.Go(func() {
					msg, ok := inportChan.Receive(ctx)
					if !ok {
						return
					}
					mu.Lock()
					names = append(names, inportName)
					fields = append(fields, msg)
					mu.Unlock()
				})
			}

			wg.Wait()

			if !outport.Send(ctx, runtime.NewStructMsg(names, fields)) {
				return
			}
		}
	}
}
