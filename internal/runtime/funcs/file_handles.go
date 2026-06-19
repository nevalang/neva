package funcs

import (
	"context"
	"errors"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

func fileHandleID(msg runtime.Msg) (int64, error) {
	idMsg, isIntMsg := msg.(runtime.IntMsg)
	if !isIntMsg {
		return 0, errors.New("file handle must be int")
	}

	return idMsg.Int(), nil
}

// sendFileHandle stores file, sends its handle, and closes it if delivery fails.
func sendFileHandle(
	ctx context.Context,
	handles *runtime.FileHandles,
	file *os.File,
	resOut runtime.SingleOutport,
	errOut runtime.SingleOutport,
) bool {
	handleID := handles.Add(file)
	if resOut.Send(ctx, runtime.NewIntMsg(handleID)) {
		return true
	}

	if err := handles.Close(handleID); err != nil {
		errOut.Send(ctx, errFromErr(err))
	}
	return false
}
