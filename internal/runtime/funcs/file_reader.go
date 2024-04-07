package funcs

import (
	"context"
	"io"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type fileReader struct{}

func (c fileReader) Create(rio runtime.FuncIO, msg runtime.Msg) (func(ctx context.Context), error) {
	filename, err := rio.In.Port("filename")
	if err != nil {
		return nil, err
	}

	dataPort, err := rio.Out.Port("data")
	if err != nil {
		return nil, err
	}

	errPort, err := rio.Out.Port("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		case name := <-filename:
			f, err := os.Open(name.Str())
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errPort <- runtime.NewStrMsg(err.Error()):
					return
				}
			}
			defer f.Close()
			data, err := io.ReadAll(f)
			if err != nil {
				select {
				case <-ctx.Done():
				case errPort <- runtime.NewStrMsg(err.Error()):
				}
			}
			select {
			case <-ctx.Done():
			case dataPort <- runtime.NewStrMsg(string(data)):
			}
		}
	}, nil
}
