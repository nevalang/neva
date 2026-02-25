package funcs

import (
	"fmt"

	"github.com/nevalang/neva/internal/runtime"
)

func errFromErr(err error) runtime.StructMsg {
	return runtime.NewStructMsg([]runtime.StructField{
		runtime.NewStructField("text", runtime.NewStringMsg(err.Error())),
	})
}

func errFromString(s string) runtime.StructMsg {
	return runtime.NewStructMsg([]runtime.StructField{
		runtime.NewStructField("text", runtime.NewStringMsg(s)),
	})
}

const (
	streamTagOpen  = "Open"
	streamTagData  = "Data"
	streamTagClose = "Close"
)

func streamOpen() runtime.UnionMsg {
	return runtime.NewUnionMsg(streamTagOpen, nil)
}

func streamData(data runtime.Msg) runtime.UnionMsg {
	return runtime.NewUnionMsg(streamTagData, data)
}

func streamClose() runtime.UnionMsg {
	return runtime.NewUnionMsg(streamTagClose, nil)
}

func streamUnion(msg runtime.Msg) runtime.UnionMsg {
	defer func() {
		if recover() != nil {
			panic(fmt.Sprintf("runtime: expected stream union message, got %T", msg))
		}
	}()

	return msg.Union()
}

func isStreamOpen(msg runtime.Msg) bool {
	return streamUnion(msg).Tag() == streamTagOpen
}

func isStreamData(msg runtime.Msg) bool {
	return streamUnion(msg).Tag() == streamTagData
}

func isStreamClose(msg runtime.Msg) bool {
	return streamUnion(msg).Tag() == streamTagClose
}

func streamDataValue(msg runtime.Msg) runtime.Msg {
	u := streamUnion(msg)
	if u.Tag() != streamTagData {
		panic("runtime: expected stream Data message")
	}
	return u.Data()
}

func emptyStruct() runtime.StructMsg {
	return runtime.NewStructMsg(nil)
}
