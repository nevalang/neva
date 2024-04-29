package funcs

import "github.com/nevalang/neva/internal/runtime"

func errorFromString(s string) runtime.Msg {
	return runtime.NewMapMsg(map[string]runtime.Msg{
		"text": runtime.NewStrMsg(s),
	})
}

func streamItem(data runtime.Msg, idx int64, last bool) runtime.Msg {
	return runtime.NewMapMsg(map[string]runtime.Msg{
		"data": data,
		"idx":  runtime.NewIntMsg(idx),
		"last": runtime.NewBoolMsg(last),
	})
}
