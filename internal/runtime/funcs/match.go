package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type match struct{}

func (match) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	ifIn, err := io.In.Array("if")
	if err != nil {
		return nil, err
	}

	thenOut, err := io.In.Array("then")
	if err != nil {
		return nil, err
	}

	if ifIn.Len() != thenOut.Len() {
		return nil, errors.New("number of 'if' inports must match number of 'then' outports")
	}

	elseIn, err := io.In.Single("else")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			ifMsgs := make([]runtime.Msg, ifIn.Len())
			if !ifIn.Receive(ctx, func(idx int, msg runtime.Msg) bool {
				ifMsgs[idx] = msg
				return true
			}) {
				return
			}

			thenMsgs := make([]runtime.Msg, thenOut.Len())
			if !thenOut.Receive(ctx, func(idx int, msg runtime.Msg) bool {
				thenMsgs[idx] = msg
				return true
			}) {
				return
			}

			elseInMsg, ok := elseIn.Receive(ctx)
			if !ok {
				return
			}

			resMsg := elseInMsg
			for i, ifMsg := range ifMsgs {
				if dataMsg == ifMsg {
					resMsg = thenMsgs[i]
					break
				}
			}

			if !resOut.Send(ctx, resMsg) {
				return
			}
		}
	}, nil
}
