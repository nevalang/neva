package funcs

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type parseInt struct{}

func (p parseInt) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	dataIn, err := io.In.Single("data")
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
			str, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			parsedNum, err := parse(str.Str())
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

func parse(str string) (runtime.Msg, error) {
	v, err := strconv.Atoi(str)
	if err != nil {
		return nil, errors.New(strings.TrimPrefix(err.Error(), "strconv.Atoi: "))
	}
	return runtime.NewIntMsg(int64(v)), nil
}
