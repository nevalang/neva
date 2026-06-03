package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// BenchmarkTernaryElseHotpath measures ternary selector cost on else path.
func BenchmarkTernaryElseHotpath(b *testing.B) {
	runtimeIO, ifIn, thenIn, elseIn, resultOutput := benchNewTernaryRuntimeIO()
	var zeroConfig runtime.Msg
	handler, err := ternarySelector{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	cancel, done := startHandler(context.Background(), handler)
	defer func() {
		cancel()
		<-done
	}()

	ifMsg := runtime.NewBoolMsg(false)
	thenMsg := runtime.NewIntMsg(42)
	elseMsg := runtime.NewIntMsg(0)

	b.ResetTimer()
	for range b.N {
		ifIn <- runtime.OrderedMsg{Msg: ifMsg}
		thenIn <- runtime.OrderedMsg{Msg: thenMsg}
		elseIn <- runtime.OrderedMsg{Msg: elseMsg}
		<-resultOutput
	}
}
