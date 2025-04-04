package funcs

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type parseFloat struct{}

func (p parseFloat) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
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

			bitsMsg, ok := bitsIn.Receive(ctx)
			if !ok {
				return
			}

			parsedNum, err := p.stringToRuntimeFloat(dataMsg, bitsMsg)
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

func (p parseFloat) stringToRuntimeFloat(
	data runtime.Msg,
	bits runtime.Msg,
) (runtime.Msg, error) {
	v, err := strconv.ParseFloat(data.Str(), int(bits.Int()))
	if err != nil {
		return nil, errors.New(strings.TrimPrefix(err.Error(), "strconv.ParseFloat: "))
	}
	return runtime.NewFloatMsg(float64(v)), nil
}
