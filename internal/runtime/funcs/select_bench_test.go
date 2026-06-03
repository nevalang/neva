package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// BenchmarkSelectHotpath measures select/then fan-in for two-slot case.
func BenchmarkSelectHotpath(b *testing.B) {
	runtimeIO, ifInputs, thenInputs, resultOutput := benchNewSelectRuntimeIO(2)
	var zeroConfig runtime.Msg
	handler, err := selector{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	cancel, done := startHandler(context.Background(), handler)
	defer func() {
		cancel()
		<-done
	}()

	ifMsg := runtime.NewBoolMsg(true)
	then0 := runtime.NewIntMsg(1)
	then1 := runtime.NewIntMsg(2)

	b.ResetTimer()
	for range b.N {
		ifInputs[0] <- runtime.OrderedMsg{Msg: ifMsg}
		thenInputs[0] <- runtime.OrderedMsg{Msg: then0}
		thenInputs[1] <- runtime.OrderedMsg{Msg: then1}
		<-resultOutput
	}
}
