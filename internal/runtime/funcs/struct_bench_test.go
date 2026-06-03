package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// BenchmarkStructHotpath measures building a small struct from three inputs.
func BenchmarkStructHotpath(b *testing.B) {
	runtimeIO, inputs, resultOutput := benchNewStructRuntimeIO([]string{"a", "b", "c"})
	var zeroConfig runtime.Msg
	handler, err := structBuilder{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	cancel, done := startHandler(context.Background(), handler)
	defer func() {
		cancel()
		<-done
	}()

	msgA := runtime.NewIntMsg(1)
	msgB := runtime.NewIntMsg(2)
	msgC := runtime.NewIntMsg(3)

	b.ResetTimer()
	for range b.N {
		inputs["a"] <- runtime.OrderedMsg{Msg: msgA}
		inputs["b"] <- runtime.OrderedMsg{Msg: msgB}
		inputs["c"] <- runtime.OrderedMsg{Msg: msgC}
		<-resultOutput
	}
}
