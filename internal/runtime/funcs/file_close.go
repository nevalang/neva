package funcs

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

type fileClose struct {
	handles *runtime.FileHandles
}

func (c fileClose) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	fileIn, err := rio.In.Single("file")
	if err != nil {
		return nil, fmt.Errorf("resolve file inport: %w", err)
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
			fileMsg, received := fileIn.Receive(ctx)
			if !received {
				return
			}

			if !c.handleFileMessage(ctx, fileMsg, resOut, errOut) {
				return
			}
		}
	}, nil
}

func (c fileClose) handleFileMessage(
	ctx context.Context,
	fileMsg runtime.OrderedMsg,
	resOut runtime.SingleOutport,
	errOut runtime.SingleOutport,
) bool {
	handleID, err := fileHandleID(fileMsg.Msg)
	if err != nil {
		return errOut.Send(ctx, errFromErr(err))
	}

	if err := c.handles.Close(handleID); err != nil {
		return errOut.Send(ctx, errFromErr(err))
	}

	return resOut.Send(ctx, emptyStruct())
}
