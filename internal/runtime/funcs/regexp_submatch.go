package funcs

import (
	"context"
	"fmt"
	"regexp"

	"github.com/nevalang/neva/internal/runtime"
)

type regexpSubmatch struct{}

func (r regexpSubmatch) Create(io runtime.FuncIO, cfgMsg runtime.Msg) (func(ctx context.Context), error) {
	regexpIn, err := io.In.SingleInport("regexp")
	if err != nil {
		return nil, err
	}

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
		var (
			regexpMsg runtime.Msg
			dataMsg   runtime.Msg
		)

		for {
			select {
			case <-ctx.Done():
				return
			case regexpMsg = <-regexpIn:
			}

			regex, err := regexp.Compile(regexpMsg.Str())
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errOut <- runtime.NewStrMsg(err.Error()):
					continue
				}
			}

			select {
			case <-ctx.Done():
				return
			case dataMsg = <-dataIn:
			}

			res := regex.FindStringSubmatch(fmt.Sprint(dataMsg))

			select {
			case <-ctx.Done():
			case resOut <- wrap(res):
			}
		}
	}, nil
}

func wrap(ss []string) runtime.Msg {
	msgs := make([]runtime.Msg, 0, len(ss))
	for _, s := range ss {
		msgs = append(msgs, runtime.NewStrMsg(s))
	}
	return runtime.NewListMsg(msgs...)
}
