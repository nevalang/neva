package funcs

import (
	"context"
	"errors"
	"os"

	"github.com/nevalang/neva/internal/runtime"
)

// parseFileHandleID extracts the opaque runtime file-handle ID from a Neva msg.
//
// File externs receive handles as Neva int values. Keeping the conversion in
// one helper gives every handle-based extern the same type check and error text.
func parseFileHandleID(msg runtime.Msg) (int64, error) {
	idMsg, isIntMsg := msg.(runtime.IntMsg)
	if !isIntMsg {
		return 0, errors.New("file handle must be int")
	}

	return idMsg.Int(), nil
}

// storeAndSendFileHandle registers an opened file and sends its handle to Neva.
//
// file_open and file_create both transfer the opened *os.File into
// runtime.FileHandles, then expose only the opaque handle ID downstream. If
// runtime cancellation prevents sending the handle, no Neva node can close that
// new handle, so this helper closes and removes it immediately to avoid leaking
// the process file.
func storeAndSendFileHandle(
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
