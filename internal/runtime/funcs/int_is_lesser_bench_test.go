package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// BenchmarkIntIsLesserHotpath measures steady-state cost of `int_is_lesser`.
func BenchmarkIntIsLesserHotpath(b *testing.B) {
	runtimeIO, leftInput, rightInput, resultOutput := benchNewBinaryRuntimeIO()
	var zeroConfig runtime.Msg
	handler, err := intIsLesser{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	cancel, done := startHandler(context.Background(), handler)
	defer func() {
		cancel()
		<-done
	}()

	left := runtime.NewIntMsg(7)
	right := runtime.NewIntMsg(42)

	b.ResetTimer()
	for range b.N {
		leftInput <- runtime.OrderedMsg{Msg: left}
		rightInput <- runtime.OrderedMsg{Msg: right}
		<-resultOutput
	}
}
