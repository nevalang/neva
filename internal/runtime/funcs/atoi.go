package funcs

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type atoi struct{}

func (a atoi) Create(io runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
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

			parsedNum, err := a.stringToRuntimeInt(str.Str())
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

func (a atoi) stringToRuntimeInt(str string) (runtime.Msg, error) {
	v, err := strconv.Atoi(str) // equivalent to ParseInt(s, 10, 0)
	if err != nil {
		return nil, errors.New(strings.TrimPrefix(err.Error(), "strconv.Atoi: "))
	}
	return runtime.NewIntMsg(int64(v)), nil
}
