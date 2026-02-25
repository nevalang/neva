package funcs

import "github.com/nevalang/neva/internal/runtime"

func errFromErr(err error) runtime.UnionMsg {
	return errFromString(err.Error())
}

func errFromString(s string) runtime.UnionMsg {
	return runtime.NewUnionMsg("Text", runtime.NewStringMsg(s))
}

func streamItem(data runtime.Msg, idx int64, last bool) runtime.StructMsg {
	return runtime.NewStructMsg([]runtime.StructField{
		runtime.NewStructField("data", data),
		runtime.NewStructField("idx", runtime.NewIntMsg(idx)),
		runtime.NewStructField("last", runtime.NewBoolMsg(last)),
	})
}

func emptyStruct() runtime.StructMsg {
	return runtime.NewStructMsg(nil)
}
