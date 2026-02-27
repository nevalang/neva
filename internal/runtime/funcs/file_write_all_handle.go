package funcs

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
)

type fileWriteAllHandle struct {
	handles *fileHandleStore
}

func (c fileWriteAllHandle) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	fileIn, err := rio.In.Single("file")
	if err != nil {
		return nil, err
	}

	dataIn, err := rio.In.Single("data")
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
			fileMsg, ok := fileIn.Receive(ctx)
			if !ok {
				return
			}

			dataMsg, ok := dataIn.Receive(ctx)
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

			if _, err := file.Write(dataMsg.Bytes()); err != nil {
				if !errOut.Send(ctx, errFromErr(err)) {
					return
				}
				continue
			}

			if !resOut.Send(ctx, runtime.NewIntMsg(id)) {
				return
			}
		}
	}, nil
}
