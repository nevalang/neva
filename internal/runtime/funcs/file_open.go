package funcs

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type fileOpen struct {
	handles *fileHandleStore
}

func (c fileOpen) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	filenameIn, err := rio.In.Single("filename")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	errOut, err := rio.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			nameMsg, ok := filenameIn.Receive(ctx)
			if !ok {
				return
			}

			file, err := os.Open(nameMsg.Str())
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			id := c.handles.Add(file)
			if !resOut.Send(ctx, runtime.NewIntMsg(id)) {
				_ = c.handles.Close(id)
				return
			}
		}
	}, nil
}
