package funcs

import (
	"context"
	"fmt"
	"io"

	"github.com/nevalang/neva/internal/runtime"
)

type fileReadAllHandle struct {
	handles *fileHandleStore
}

func (c fileReadAllHandle) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	fileIn, err := rio.In.Single("file")
	if err != nil {
		return nil, fmt.Errorf("resolve file inport: %w", err)
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, fmt.Errorf("resolve res outport: %w", err)
	}

	handleOut, err := rio.Out.Single("handle")
	if err != nil {
		return nil, fmt.Errorf("resolve handle outport: %w", err)
	}

	errOut, err := rio.Out.Single("err")
	if err != nil {
		return nil, fmt.Errorf("resolve err outport: %w", err)
	}

	return func(ctx context.Context) {
		for {
			fileMsg, received := fileIn.Receive(ctx)
			if !received {
				return
			}

			if !c.handleFileMessage(ctx, fileMsg, resOut, handleOut, errOut) {
				return
			}
		}
	}, nil
}

func (c fileReadAllHandle) handleFileMessage(
	ctx context.Context,
	fileMsg runtime.OrderedMsg,
	resOut runtime.SingleOutport,
	handleOut runtime.SingleOutport,
	errOut runtime.SingleOutport,
) bool {
	handleID, err := fileHandleID(fileMsg.Msg)
	if err != nil {
		return errOut.Send(ctx, errFromErr(err))
	}

	file, err := c.handles.Get(handleID)
	if err != nil {
		return errOut.Send(ctx, errFromErr(err))
	}

	data, err := io.ReadAll(file)
	if err != nil {
		if !handleOut.Send(ctx, runtime.NewIntMsg(handleID)) {
			return false
		}
		return errOut.Send(ctx, errFromErr(err))
	}

	if !resOut.Send(ctx, runtime.NewBytesMsg(data)) {
		return false
	}

	return handleOut.Send(ctx, runtime.NewIntMsg(handleID))
}
