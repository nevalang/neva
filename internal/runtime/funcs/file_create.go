package funcs

import (
	"context"
	"fmt"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

type fileCreate struct {
	handles *runtime.FileHandles
}

func (c fileCreate) Create(rio runtime.IO, _ runtime.Msg) (func(ctx context.Context), error) {
	filenameIn, err := rio.In.Single("filename")
	if err != nil {
		return nil, fmt.Errorf("resolve filename inport: %w", err)
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
			nameMsg, received := filenameIn.Receive(ctx)
			if !received {
				return
			}

			if !c.handleFileMessage(ctx, nameMsg, resOut, errOut) {
				return
			}
		}
	}, nil
}

func (c fileCreate) handleFileMessage(
	ctx context.Context,
	nameMsg runtime.OrderedMsg,
	resOut runtime.SingleOutport,
	errOut runtime.SingleOutport,
) bool {
	// #nosec G304 -- filename is user-controlled by design.
	file, err := os.OpenFile(nameMsg.Str(), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o644)
	if err != nil {
		return errOut.Send(ctx, errFromErr(err))
	}

	handleID := c.handles.Add(file)
	if resOut.Send(ctx, runtime.NewIntMsg(handleID)) {
		return true
	}

	if err := c.handles.Close(handleID); err != nil {
		panic(err)
	}
	return false
}
