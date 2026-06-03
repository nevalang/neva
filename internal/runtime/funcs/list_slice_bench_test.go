package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// BenchmarkListSliceHotpath measures list slicing with fixed bounds.
func BenchmarkListSliceHotpath(b *testing.B) {
	runtimeIO, dataIn, fromIn, toIn, resultOutput := benchNewListSliceRuntimeIO()
	var zeroConfig runtime.Msg
	handler, err := listSlice{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	cancel, done := startHandler(context.Background(), handler)
	defer func() {
		cancel()
		<-done
	}()

	data := runtime.NewListMsg([]runtime.Msg{
		runtime.NewIntMsg(1),
		runtime.NewIntMsg(2),
		runtime.NewIntMsg(3),
		runtime.NewIntMsg(4),
	})
	from := runtime.NewIntMsg(1)
	to := runtime.NewIntMsg(3)

	b.ResetTimer()
	for range b.N {
		dataIn <- runtime.OrderedMsg{Msg: data}
		fromIn <- runtime.OrderedMsg{Msg: from}
		toIn <- runtime.OrderedMsg{Msg: to}
		<-resultOutput
	}
}
