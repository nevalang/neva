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
	dataIn, err := io.In.SingleInport("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.SingleOutport("res")
	if err != nil {
		return nil, err
	}

	errOut, err := io.Out.SingleOutport("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var str runtime.Msg

		for {
			select {
			case <-ctx.Done():
				return
			case str = <-dataIn:
			}

			parsedNum, err := parse(str.Str())
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errOut <- runtime.NewStrMsg(err.Error()):
				}
				continue
			}

			select {
			case <-ctx.Done():
				return
			case resOut <- parsedNum:
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
