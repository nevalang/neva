package funcs

import (
	"context"
	"regexp"

	"github.com/nevalang/neva/internal/runtime"
)

type regexpSubmatcher struct{}

func (r regexpSubmatcher) Create(io runtime.FuncIO, cfgMsg runtime.Msg) (func(ctx context.Context), error) {
	regexpIn, err := io.In.Port("regexp")
	if err != nil {
		return nil, err
	}

	dataIn, err := io.In.Port("data")
	if err != nil {
		return nil, err
	}

	resOut, err := io.Out.Port("res")
	if err != nil {
		return nil, err
	}

	errOut, err := io.Out.Port("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case regexpMsg := <-regexpIn:
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
				case dataMsg := <-dataIn:
					resOut <- wrap(
						regex.FindStringSubmatch(
							dataMsg.String(),
						),
					)
				}
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
