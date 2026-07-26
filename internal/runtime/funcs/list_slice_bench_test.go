package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

// BenchmarkListSliceHotpath measures list slicing with fixed bounds.
func BenchmarkListSliceHotpath(b *testing.B) {
	runtimeIO, dataIn, fromIn, toIn, resultOutput := benchNewListSliceRuntimeIO()
	var zeroConfig messages.Msg
	handler, err := listSlice{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	cancel, done := startHandler(context.Background(), handler)
	defer func() {
		cancel()
		<-done
	}()

	data := messages.NewListIntMsg([]int64{1, 2, 3, 4})
	from := messages.NewIntMsg(1)
	to := messages.NewIntMsg(3)

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		dataIn <- runtime.OrderedMsg{Msg: data}
		fromIn <- runtime.OrderedMsg{Msg: from}
		toIn <- runtime.OrderedMsg{Msg: to}
		<-resultOutput
	}
}
