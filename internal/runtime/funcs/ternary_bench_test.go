package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

// BenchmarkTernaryElseHotpath measures ternary selector cost on else path.
func BenchmarkTernaryElseHotpath(b *testing.B) {
	runtimeIO, ifIn, thenIn, elseIn, resultOutput := benchNewTernaryRuntimeIO()
	var zeroConfig messages.Msg
	handler, err := ternarySelector{}.Create(runtimeIO, zeroConfig)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	cancel, done := startHandler(context.Background(), handler)
	defer func() {
		cancel()
		<-done
	}()

	ifMsg := messages.NewBoolMsg(false)
	thenMsg := messages.NewIntMsg(42)
	elseMsg := messages.NewIntMsg(0)

	b.ResetTimer()
	for range b.N {
		ifIn <- runtime.OrderedMsg{Msg: ifMsg}
		thenIn <- runtime.OrderedMsg{Msg: thenMsg}
		elseIn <- runtime.OrderedMsg{Msg: elseMsg}
		<-resultOutput
	}
}
