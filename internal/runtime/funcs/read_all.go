package funcs

import (
	"context"
	"io"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type readAll struct{}

func (c readAll) Create(rio runtime.FuncIO, _ runtime.Msg) (func(ctx context.Context), error) {
	filename, err := rio.In.Single("filename")
	if err != nil {
		return nil, err
	}

	dataPort, err := rio.Out.Single("data")
	if err != nil {
		return nil, err
	}

	errPort, err := rio.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			name, ok := filename.Receive(ctx)
			if !ok {
				return
			}

			f, err := os.Open(name.Str())
			if err != nil {
				if !errPort.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			data, err := io.ReadAll(f)
			if err != nil {
				if !errPort.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if err := f.Close(); err != nil {
				if !errPort.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !dataPort.Send(ctx, runtime.NewStrMsg(string(data))) {
				return
			}
		}
	}, nil
}
