package runtime

import (
	"context"
	"testing"
)

// BenchmarkArrayInportSelectTwoSlots measures the two-slot fast-path in Select:
// one blocking receive + opportunistic competitor probe + index-based ordering.
func BenchmarkArrayInportSelectTwoSlots(b *testing.B) {
	ch0 := make(chan OrderedMsg, b.N+1)
	ch1 := make(chan OrderedMsg, b.N+1)
	in := newTestArrayInport(ch0, ch1)
	ctx := context.Background()

	for i := range b.N {
		ch0 <- OrderedMsg{Msg: NewIntMsg(1), index: uint64(i*2 + 2)}
		ch1 <- OrderedMsg{Msg: NewIntMsg(2), index: uint64(i*2 + 1)}

		if _, ok := in.Select(ctx); !ok {
			b.Fatal("unexpected canceled select")
		}
		if _, ok := in.Select(ctx); !ok {
			b.Fatal("unexpected canceled select")
		}
	}
}
