package funcs

import (
	"errors"

	"github.com/nevalang/neva/internal/runtime"
)

func fileHandleID(msg runtime.Msg) (int64, error) {
	idMsg, isIntMsg := msg.(runtime.IntMsg)
	if !isIntMsg {
		return 0, errors.New("file handle must be int")
	}

	return idMsg.Int(), nil
}
