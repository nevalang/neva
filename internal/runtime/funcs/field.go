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

	msgIn, err := io.In.SingleInport("msg")
	if err != nil {
		return nil, err
	}

	msgOut, err := io.Out.SingleOutport("msg")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		var msg runtime.Msg

		for {
			select {
			case <-ctx.Done():
				return
			case msg = <-msgIn:
			}

			select {
			case <-ctx.Done():
				return
			case msgOut <- s.recursive(msg, pathStrings):
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
