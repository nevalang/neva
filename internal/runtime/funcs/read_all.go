package funcs

import (
	"context"
	"io"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type readAll struct{}

func (c readAll) Create(rio runtime.FuncIO, msg runtime.Msg) (func(ctx context.Context), error) {
	filename, err := rio.In.SingleInport("filename")
	if err != nil {
		return nil, err
	}

	dataPort, err := rio.Out.SingleOutport("data")
	if err != nil {
		return nil, err
	}

	errPort, err := rio.Out.SingleOutport("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var name runtime.Msg
			select {
			case <-ctx.Done():
				return
			case name = <-filename:
			}

			f, err := os.Open(name.Str())
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errPort <- errorFromString(err.Error()):
					continue
				}
			}

			data, err := io.ReadAll(f)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errPort <- errorFromString(err.Error()):
					continue
				}
			}

			if err := f.Close(); err != nil {
				select {
				case <-ctx.Done():
					return
				case errPort <- errorFromString(err.Error()):
					continue
				}
			}

			select {
			case <-ctx.Done():
				return
			case dataPort <- runtime.NewStrMsg(string(data)):
			}
		}
	}, nil
}
