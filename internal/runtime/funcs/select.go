package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type selector struct{}

func (selector) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	ifArrIn, err := io.In.Array("if")
	if err != nil {
		return nil, err
	}

	thenArrIn, err := io.In.Array("then")
	if err != nil {
		return nil, err
	}

	if ifArrIn.Len() != thenArrIn.Len() {
		return nil, errors.New("number of 'if' inports must match number of 'then' outports")
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			ifMsg, ok := ifArrIn.Select(ctx)
			if !ok {
				return
			}

			then := make([]runtime.Msg, thenArrIn.Len())
			if !thenArrIn.ReceiveAll(ctx, func(idx int, msg runtime.Msg) bool {
				then[idx] = msg
				return true
			}) {
				return
			}

			if !resOut.Send(ctx, then[ifMsg.SlotIdx]) {
				return
			}
		}
	}, nil
}
