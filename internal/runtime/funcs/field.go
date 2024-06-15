package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type readStructField struct{}

func (s readStructField) Create(io runtime.FuncIO, cfg runtime.Msg) (func(ctx context.Context), error) {
	path := cfg.List()
	if len(path) == 0 {
		return nil, errors.New("field path cannot be empty")
	}

	pathStrings := make([]string, 0, len(path))
	for _, el := range path {
		pathStrings = append(pathStrings, el.Str())
	}

	msgIn, err := io.In.Single("msg")
	if err != nil {
		return nil, err
	}

	msgOut, err := io.Out.SingleOutport("msg")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			msg, ok := msgIn.Receive(ctx)
			if !ok {
				return
			}

			if !msgOut.Send(
				ctx,
				s.recursive(msg, pathStrings),
			) {
				return
			}
		}
	}, nil
}

func (readStructField) recursive(m runtime.Msg, path []string) runtime.Msg {
	for len(path) > 0 {
		m = m.Map()[path[0]]
		path = path[1:]
	}
	return m
}
