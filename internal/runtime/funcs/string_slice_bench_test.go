package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

// BenchmarkStringSliceHotpath measures string slicing with fixed bounds.
func BenchmarkStringSliceHotpath(b *testing.B) {
	runtimeIO, dataIn, fromIn, toIn, resultOutput := benchNewStringSliceRuntimeIO()
	var zeroConfig messages.Msg
	handler, err := stringSlice{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	cancel, done := startHandler(context.Background(), handler)
	defer func() {
		cancel()
		<-done
	}()

	data := messages.NewStringMsg("abcd")
	from := messages.NewIntMsg(1)
	to := messages.NewIntMsg(3)

	b.ResetTimer()
	for range b.N {
		dataIn <- runtime.OrderedMsg{Msg: data}
		fromIn <- runtime.OrderedMsg{Msg: from}
		toIn <- runtime.OrderedMsg{Msg: to}
		<-resultOutput
	}
}
