package funcs

import "github.com/nevalang/neva/internal/runtime"

func errFromErr(err error) runtime.StructMsg {
	return runtime.NewStructMsg(
		[]string{"text"},
		[]runtime.Msg{runtime.NewStrMsg(err.Error())},
	)
}

func errFromString(s string) runtime.StructMsg {
	return runtime.NewStructMsg(
		[]string{"text"},
		[]runtime.Msg{runtime.NewStrMsg(s)},
	)
}

func streamItem(data runtime.Msg, idx int64, last bool) runtime.StructMsg {
	return runtime.NewStructMsg(
		[]string{"data", "idx", "last"},
		[]runtime.Msg{data, runtime.NewIntMsg(idx), runtime.NewBoolMsg(last)},
	)
}
