package funcs

import (
	"context"
	"io"

	"github.com/nevalang/neva/internal/runtime"
)

type fileReadAllHandle struct {
	handles *fileHandleStore
}

func (c fileReadAllHandle) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	fileIn, err := rio.In.Single("file")
	if err != nil {
		return nil, err
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, err
	}

	handleOut, err := rio.Out.Single("handle")
	if err != nil {
		return nil, err
	}

	errOut, err := rio.Out.Single("err")
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		for {
			fileMsg, ok := fileIn.Receive(ctx)
			if !ok {
				return
			}

			id, err := fileHandleID(fileMsg)
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			file, err := c.handles.Get(id)
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			data, err := io.ReadAll(file)
			if err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, runtime.NewStringMsg(string(data))) {
				return
			}
			if !handleOut.Send(ctx, runtime.NewIntMsg(id)) {
				return
			}
		}
	}, nil
}
