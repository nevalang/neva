package funcs

import "github.com/nevalang/neva/internal/runtime"

// tryToUnboxIfUnion unwraps union payload when it is present.
// For no-data unions it returns original union message unchanged.
//
//nolint:ireturn // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func tryToUnboxIfUnion(msg runtime.Msg) runtime.Msg {
	unionMsg, ok := runtime.AsUnion(msg)
	if !ok {
		return msg
	}

	if unionMsg.Data() == nil {
		return msg
	}

	return unionMsg.Data()
}
