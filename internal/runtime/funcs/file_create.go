package funcs

import (
	"context"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type fileCreate struct {
	handles *fileHandleStore
}

func (c fileCreate) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
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

			// #nosec G304 -- filename is user-controlled by design.
			file, err := os.OpenFile(nameMsg.Str(), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o644)
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
