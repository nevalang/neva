package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type structField struct{}

//nolint:varnamelen // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (s structField) Create(io runtime.IO, cfg runtime.Msg) (func(ctx context.Context), error) {
	typedPath, ok := runtime.AsListStrings(cfg.List())
	if !ok {
		return nil, errors.New("field config must be list<string>")
	}
	pathStrings := append([]string(nil), typedPath...)

	if len(pathStrings) == 0 {
		return nil, errors.New("field path cannot be empty")
	}

	dataIn, err := io.In.Single("data")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	resOut, err := io.Out.Single("res")
	if err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			dataMsg, ok := dataIn.Receive(ctx)
			if !ok {
				return
			}

			if !resOut.Send(ctx, s.selector(dataMsg, pathStrings)) {
				return
			}
		}
	}, nil
}

//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (structField) selector(m runtime.Msg, path []string) runtime.Msg {
	for len(path) > 0 {
		m = m.Struct().Get(path[0])
		path = path[1:]
	}
	return m
}
