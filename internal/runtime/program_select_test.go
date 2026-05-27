package runtime

import (
	"context"
	"testing"
)

func newTestArrayInport(chans ...<-chan OrderedMsg) *ArrayInport {
	return NewArrayInport(
		NewTracer(),
		chans,
		PortAddr{Path: "test/in", Port: "data"},
		NoEffectInterceptor{},
	)
}

func TestArrayInportSelectSingleSlot(t *testing.T) {
	t.Parallel()

	ch := make(chan OrderedMsg, 1)
	ch <- OrderedMsg{Msg: NewIntMsg(7), index: 10}
	in := newTestArrayInport(ch)

	selected, ok := in.Select(context.Background())
	if !ok {
		t.Fatal("expected message")
	}
	if selected.SlotIdx != 0 {
		t.Fatalf("unexpected slot index: %d", selected.SlotIdx)
	}
	if selected.OrderedMsg.Int() != 7 {
		t.Fatalf("unexpected payload: %d", selected.OrderedMsg.Int())
	}
}

func TestArrayInportSelectTwoSlotsOrdersByIndex(t *testing.T) {
	t.Parallel()

	ch0 := make(chan OrderedMsg, 2)
	ch1 := make(chan OrderedMsg, 2)
	ch0 <- OrderedMsg{Msg: NewStringMsg("new"), index: 200}
	ch1 <- OrderedMsg{Msg: NewStringMsg("old"), index: 100}
	in := newTestArrayInport(ch0, ch1)

	first, ok := in.Select(context.Background())
	if !ok {
		t.Fatal("expected first message")
	}
	if first.OrderedMsg.Str() != "old" {
		t.Fatalf("unexpected first payload: %q", first.OrderedMsg.Str())
	}

	second, ok := in.Select(context.Background())
	if !ok {
		t.Fatal("expected second message")
	}
	if second.OrderedMsg.Str() != "new" {
		t.Fatalf("unexpected second payload: %q", second.OrderedMsg.Str())
	}
}

func TestArrayInportSelectContextCancel(t *testing.T) {
	t.Parallel()

	ch0 := make(chan OrderedMsg)
	ch1 := make(chan OrderedMsg)
	in := newTestArrayInport(ch0, ch1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, ok := in.Select(ctx)
	if ok {
		t.Fatal("expected select to stop on canceled context")
	}
}

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
