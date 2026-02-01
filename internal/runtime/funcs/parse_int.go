package funcs

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type parseInt struct{}

func (p parseInt) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
	if err != nil {
		return nil, err
	}

	baseIn, err := io.In.Single("base")
	if err != nil {
		return nil, err
	}

	bitsIn, err := io.In.Single("bits")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		return nil, err
	}

	errOut, err := io.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			baseMsg, ok := baseIn.Receive(ctx)
			if !ok {
				return
			}

			bitsMsg, ok := bitsIn.Receive(ctx)
			if !ok {
				return
			}

			parsedNum, err := p.stringToRuntimeInt(dataMsg, baseMsg, bitsMsg)
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, parsedNum) {
				return
			}
		}
	}, nil
}

func (p parseInt) stringToRuntimeInt(
	data runtime.Msg,
	base runtime.Msg,
	bits runtime.Msg,
) (runtime.Msg, error) {
	v, err := strconv.ParseInt(
		data.Str(),
		int(base.Int()),
		int(bits.Int()),
	)
	if err != nil {
		return nil, errors.New(strings.TrimPrefix(err.Error(), "strconv.Atoi: "))
	}
	return runtime.NewIntMsg(v), nil
}
