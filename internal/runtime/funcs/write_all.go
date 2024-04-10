package funcs

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type writeAll struct{}

func (c writeAll) Create(rio runtime.FuncIO, msg runtime.Msg) (func(ctx context.Context), error) {
	filename, err := rio.In.Port("filename")
	if err != nil {
		return nil, err
	}

	dataPort, err := rio.In.Port("data")
	if err != nil {
		return nil, err
	}

	sig, err := rio.Out.Port("sig")
	if err != nil {
		return nil, err
	}

	errPort, err := rio.Out.Port("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			var name, data runtime.Msg
			select {
			case <-ctx.Done():
				return
			case name = <-filename:
			}
			if name.Type() != runtime.StrMsgType {
				select {
				case <-ctx.Done():
					return
				case errPort <- runtime.NewMapMsg(map[string]runtime.Msg{
					"text": runtime.NewStrMsg("filename must be a string"),
				}):
					continue
				}
			}
			select {
			case <-ctx.Done():
				return
			case data = <-dataPort:
			}
			if name.Type() != runtime.StrMsgType {
				select {
				case <-ctx.Done():
					return
				case errPort <- runtime.NewMapMsg(map[string]runtime.Msg{
					"text": runtime.NewStrMsg("data must be a string"),
				}):
					continue
				}
			}
			if err := os.WriteFile(name.Str(), []byte(data.Str()), 0755); err != nil {
				select {
				case <-ctx.Done():
					return
				case errPort <- runtime.NewMapMsg(map[string]runtime.Msg{
					"text": runtime.NewStrMsg("filename must be a string"),
				}):
					continue
				}
			}
			select {
			case <-ctx.Done():
				return
			case sig <- nil:
			}
		}
	}, nil
}
