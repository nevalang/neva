package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type selector struct{}

func (selector) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	ifIn, err := io.In.Array("if")
	if err != nil {
		return nil, err
	}

	thenIn, err := io.In.Array("then")
	if err != nil {
		return nil, err
	}

	if ifIn.Len() != thenIn.Len() {
		return nil, errors.New("number of 'if' inports must match number of 'then' outports")
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			then := make([]runtime.Msg, ifIn.Len())
			if !thenIn.Receive(ctx, func(idx int, msg runtime.Msg) bool {
				then[idx] = msg
				return true
			}) {
				return
			}

			ifMsgs, ok := ifIn.Select(ctx)
			if !ok {
				return
			}

			for _, ifMsg := range ifMsgs {
				if !resOut.Send(ctx, then[ifMsg.SlotIdx]) {
					return
				}
			}
		}
	}, nil
}
