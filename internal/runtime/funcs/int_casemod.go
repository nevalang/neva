package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type intCaseMod struct{}

func (intCaseMod) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	caseIn, err := io.In.Array("case")
	if err != nil {
		return nil, err
	}

	caseOut, err := io.Out.Array("case")
	if err != nil {
		return nil, err
	}

	elseOut, err := io.Out.Single("else")
	if err != nil {
		return nil, err
	}

	if caseIn.Len() != caseOut.Len() {
		return nil, errors.New("number of 'case' inports must match number of 'case' outports")
	}

	return func(ctx context.Context) {
		for {
			data, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			matchIdx := -1
			dataInt := data.Int()

			if !caseIn.Receive(ctx, func(idx int, msg runtime.Msg) bool {
				if matchIdx == -1 && dataInt%msg.Int() == 0 {
					matchIdx = idx
				}
				return true
			}) {
				return
			}

			if matchIdx != -1 {
				if !caseOut.Send(ctx, uint8(matchIdx), data) {
					return
				}
			} else {
				if !elseOut.Send(ctx, data) {
					return
				}
			}
		}
	}, nil
}
