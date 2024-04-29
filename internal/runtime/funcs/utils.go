package funcs

import "github.com/nevalang/neva/internal/runtime"

func errorFromString(s string) runtime.Msg {
	return runtime.NewMapMsg(map[string]runtime.Msg{
		"text": runtime.NewStrMsg(s),
	})
}
