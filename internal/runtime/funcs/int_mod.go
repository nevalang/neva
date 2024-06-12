package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type intMod struct{}

func (intMod) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.SingleInport("data")
	if err != nil {
		return nil, err
	}

	caseIn, ok := io.In["case"]
	if !ok {
		return nil, errors.New("port 'case' is required")
	}

	caseOut, ok := io.Out["case"]
	if !ok {
		return nil, errors.New("port 'then' is required")
	}

	elseOut, err := io.Out.SingleOutport("else")
	if err != nil {
		return nil, err
	}

	if len(caseIn) != len(caseOut) {
		return nil, errors.New("number of 'case' inports must match number of 'then' outports")
	}

	return func(ctx context.Context) {
		var data runtime.Msg

		for {
			select {
			case <-ctx.Done():
				return
			case data = <-dataIn:
			}

			cases := make([]runtime.Msg, len(caseIn))
			for i, slot := range caseIn {
				select {
				case <-ctx.Done():
					return
				case caseMsg := <-slot:
					cases[i] = caseMsg
				}
			}

			matchIdx := -1
			dataInt := data.Int()
			for i, caseMsg := range cases {
				if dataInt%caseMsg.Int() == 0 {
					matchIdx = i
					break
				}
			}

			if matchIdx != -1 {
				select {
				case <-ctx.Done():
					return
				case caseOut[matchIdx] <- data:
					continue
				}
			}

			select {
			case <-ctx.Done():
				return
			case elseOut <- data:
			}
		}
	}, nil
}
