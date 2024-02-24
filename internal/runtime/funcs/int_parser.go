package funcs

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/nevalang/neva/internal/runtime"
)

type intParser struct{}

func (p intParser) Create(io runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	vin, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	vout, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	errout, err := io.Out.Port("err")
	if err != nil {
		return nil, err
	}

	parseFunc := func(str string) (runtime.Msg, error) {
		v, err := strconv.Atoi(str)
		if err != nil {
			return nil, errors.New(strings.TrimPrefix(err.Error(), "strconv.Atoi: "))
		}
		return runtime.NewIntMsg(int64(v)), nil
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case str := <-vin:
				v, err := parseFunc(str.Str())
				if err != nil {
					select {
					case <-ctx.Done():
						return
					case errout <- runtime.NewStrMsg(err.Error()):
					}
					continue
				}
				select {
				case <-ctx.Done():
					return
				case vout <- v:
				}
			}
		}
	}, nil
}
