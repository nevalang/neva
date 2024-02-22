package funcs

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

type structSelector struct{}

func (s structSelector) Create(io runtime.FuncIO, msg runtime.Msg) (func(ctx context.Context), error) {
	fieldPath := msg.List()
	if len(fieldPath) == 0 {
		return nil, errors.New("field path cannot be empty")
	}

	path := make([]string, 0, len(fieldPath))
	for _, el := range fieldPath {
		path = append(path, el.Str())
	}

	vin, err := io.In.Port("v")
	if err != nil {
		return nil, err
	}

	vout, err := io.Out.Port("v")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case obj := <-vin:
				select {
				case <-ctx.Done():
					return
				case vout <- s.getFieldByPath(obj, path):
				}
			}
		}
	}, nil
}

func (structSelector) getFieldByPath(msg runtime.Msg, fieldPath []string) runtime.Msg {
	for len(fieldPath) > 0 {
		msg = msg.Map()[fieldPath[0]]
		fieldPath = fieldPath[1:]
	}
	return msg
}
