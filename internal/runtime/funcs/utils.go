package funcs

import "github.com/nevalang/neva/internal/runtime"

func errFromErr(err error) runtime.MapMsg {
	return runtime.NewMapMsg(map[string]runtime.Msg{
		"text": runtime.NewStrMsg(err.Error()),
	})
}

func errFromString(s string) runtime.MapMsg {
	return runtime.NewMapMsg(map[string]runtime.Msg{
		"text": runtime.NewStrMsg(s),
	})
}

func streamItem(data runtime.Msg, idx int64, last bool) runtime.MapMsg {
	return runtime.NewMapMsg(map[string]runtime.Msg{
		"data": data,
		"idx":  runtime.NewIntMsg(idx),
		"last": runtime.NewBoolMsg(last),
	})
}
