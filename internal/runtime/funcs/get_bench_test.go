package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// BenchmarkGetDictTypedInt measures a complete typed dictionary lookup through ports.
func BenchmarkGetDictTypedInt(b *testing.B) {
	io, inputs, outputs := newIO([]string{"dict", "key"}, []string{"res", "err"})
	handler, err := (getDictValue{}).Create(io, nil)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}
	ctx, cancel := context.WithCancel(b.Context())
	defer cancel()
	go handler(ctx)

	dict := runtime.NewDictIntMsg(map[string]int64{"answer": 42})
	key := runtime.NewStringMsg("answer")
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		inputs["dict"] <- runtime.OrderedMsg{Msg: dict}
		inputs["key"] <- runtime.OrderedMsg{Msg: key}
		<-outputs["res"]
	}
}
