package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

// BenchmarkIntIsLesserHotpath measures steady-state cost of `int_is_lesser`.
func BenchmarkIntIsLesserHotpath(b *testing.B) {
	runtimeIO, leftInput, rightInput, resultOutput := benchNewBinaryRuntimeIO()
	var zeroConfig messages.Msg
	handler, err := intIsLesser{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	cancel, done := startHandler(context.Background(), handler)
	defer func() {
		cancel()
		<-done
	}()

	left := messages.NewIntMsg(7)
	right := messages.NewIntMsg(42)

	b.ResetTimer()
	for range b.N {
		leftInput <- runtime.OrderedMsg{Msg: left}
		rightInput <- runtime.OrderedMsg{Msg: right}
		<-resultOutput
	}
}
