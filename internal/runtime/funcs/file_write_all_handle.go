package funcs

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type fileWriteAllHandle struct {
	handles *runtime.FileHandles
}

func (c fileWriteAllHandle) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	fileIn, err := rio.In.Single("file")
	if err != nil {
		return nil, fmt.Errorf("resolve file inport: %w", err)
	}

	dataIn, err := rio.In.Single("data")
	if err != nil {
		return nil, fmt.Errorf("resolve data inport: %w", err)
	}

	resOut, err := rio.Out.Single("res")
	if err != nil {
		return nil, fmt.Errorf("resolve res outport: %w", err)
	}

	errOut, err := rio.Out.Single("err")
	if err != nil {
		return nil, fmt.Errorf("resolve err outport: %w", err)
	}

	return func(ctx context.Context) {
		for {
			fileMsg, dataMsg, received := receive2(ctx, fileIn, dataIn)
			if !received {
				return
			}

			if !c.handleFileMessage(ctx, fileMsg, dataMsg, resOut, errOut) {
				return
			}
		}
	}, nil
}

func (c fileWriteAllHandle) handleFileMessage(
	ctx context.Context,
	fileMsg runtime.OrderedMsg,
	dataMsg runtime.OrderedMsg,
	resOut runtime.SingleOutport,
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

	if _, err := file.Write(dataMsg.Bytes()); err != nil {
		if !resOut.Send(ctx, runtime.NewIntMsg(handleID)) {
			return false
		}
		return errOut.Send(ctx, errFromErr(err))
	}

	return resOut.Send(ctx, runtime.NewIntMsg(handleID))
}
