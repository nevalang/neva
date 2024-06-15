package funcs

import (
	"context"
	"fmt"
	"regexp"

	"github.com/nevalang/neva/internal/runtime"
)

type regexpSubmatch struct{}

func (r regexpSubmatch) Create(io runtime.FuncIO, cfgMsg runtime.Msg) (func(ctx context.Context), error) {
	regexpIn, err := io.In.Single("regexp")
	if err != nil {
		return nil, err
	}

	dataIn, err := io.In.Single("data")
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
		for {
			regexpMsg, ok := regexpIn.Receive(ctx)
			if !ok {
				return
			}

			regex, err := regexp.Compile(regexpMsg.Str())
			if err != nil {
				if !errOut.Send(ctx, runtime.NewStrMsg(err.Error())) {
					return
				}
				continue
			}

			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(
				ctx,
				stringsToList(
					regex.FindStringSubmatch(
						fmt.Sprint(dataMsg),
					),
				),
			) {
				return
			}
		}
	}, nil
}

func stringsToList(ss []string) runtime.Msg {
	msgs := make([]runtime.Msg, 0, len(ss))
	for _, s := range ss {
		msgs = append(msgs, runtime.NewStrMsg(s))
	}
	return runtime.NewListMsg(msgs)
}
