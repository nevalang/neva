package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type race struct{}

func (race) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	casesIn, err := io.In.Array("case")
	if err != nil {
		return nil, err
	}

	casesOut, err := io.Out.Array("case")
	if err != nil {
		return nil, err
	}

	if casesIn.Len() != casesOut.Len() {
		return nil, errors.New("number of 'case' inports must match number of 'case' outports")
	}

	return func(ctx context.Context) {
		for {
			msg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			caseMsg, ok := casesIn.Select(ctx)
			if !ok {
				return
			}

			if !casesOut.Send(ctx, caseMsg.SlotIdx, msg) {
				return
			}
		}
	}, nil
}
