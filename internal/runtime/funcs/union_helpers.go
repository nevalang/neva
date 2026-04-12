package funcs

import "github.com/nevalang/neva/internal/runtime"

func tryToUnboxIfUnion(msg runtime.Msg) runtime.Msg {
	if !msg.IsUnion() {
		return msg
	}

	unionMsg := msg.Union()
	if !unionMsg.HasData() {
		return msg
	}

	return unionMsg.Data()
}
