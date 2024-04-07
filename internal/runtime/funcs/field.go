package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type field struct{}

func (s field) Create(io runtime.FuncIO, fieldPathMsg runtime.Msg) (func(ctx context.Context), error) {
	fieldPath := fieldPathMsg.List()
	if len(fieldPath) == 0 {
		return nil, errors.New("field path cannot be empty")
	}

	pathStrings := make([]string, 0, len(fieldPath))
	for _, el := range fieldPath {
		pathStrings = append(pathStrings, el.Str())
	}

	msgIn, err := io.In.Port("msg")
	if err != nil {
		return nil, err
	}

	msgOut, err := io.Out.Port("msg")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgIn:
				select {
				case <-ctx.Done():
					return
				case msgOut <- s.mapLookup(msg, pathStrings):
				}
			}
		}
	}, nil
}

func (field) mapLookup(m runtime.Msg, path []string) runtime.Msg {
	for len(path) > 0 {
		m = m.Map()[path[0]]
		path = path[1:]
	}
	return m
}
