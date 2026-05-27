package funcs

import (
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

// TestListAtTypedInt verifies typed int-list fast path and negative indexing invariant.
func TestListAtTypedInt(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newIO([]string{"data", "idx"}, []string{"res", "err"})
	handler, err := (listAt{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	cancel, done := runHandler(handler)
	defer func() {
		cancel()
		<-done
	}()

	inChans["data"] <- runtime.OrderedMsg{Msg: runtime.NewListIntMsg([]int64{10, 20, 30})}
	inChans["idx"] <- runtime.OrderedMsg{Msg: runtime.NewIntMsg(-1)}

	select {
	case got := <-outChans["res"]:
		if got.Int() != 30 {
			t.Fatalf("result = %d, want 30", got.Int())
		}
	case <-time.After(time.Second):
		t.Fatal("no result")
	}
}

// TestListAtGenericFallback verifies generic list path still works for mixed/non-typed lists.
func TestListAtGenericFallback(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newIO([]string{"data", "idx"}, []string{"res", "err"})
	handler, err := (listAt{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	cancel, done := runHandler(handler)
	defer func() {
		cancel()
		<-done
	}()

	inChans["data"] <- runtime.OrderedMsg{Msg: runtime.NewListMsg([]runtime.Msg{runtime.NewStringMsg("a"), runtime.NewIntMsg(2)})}
	inChans["idx"] <- runtime.OrderedMsg{Msg: runtime.NewIntMsg(1)}

	select {
	case got := <-outChans["res"]:
		if got.Int() != 2 {
			t.Fatalf("result = %d, want 2", got.Int())
		}
	case <-time.After(time.Second):
		t.Fatal("no result")
	}
}

// TestListAtOutOfBounds verifies out-of-bounds invariant is routed to err outport.
func TestListAtOutOfBounds(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newIO([]string{"data", "idx"}, []string{"res", "err"})
	handler, err := (listAt{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	cancel, done := runHandler(handler)
	defer func() {
		cancel()
		<-done
	}()

	inChans["data"] <- runtime.OrderedMsg{Msg: runtime.NewListStringMsg([]string{"x"})}
	inChans["idx"] <- runtime.OrderedMsg{Msg: runtime.NewIntMsg(5)}

	select {
	case got := <-outChans["err"]:
		if got.Struct().Get("text").Str() != "index out of bounds" {
			t.Fatalf("error text = %q, want %q", got.Struct().Get("text").Str(), "index out of bounds")
		}
	case <-time.After(time.Second):
		t.Fatal("no error result")
	}

	select {
	case <-outChans["res"]:
		t.Fatal("unexpected success result")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestListItem(t *testing.T) {
	t.Parallel()

	item, ok := listItem([]int{1, 2, 3}, -2)
	if !ok || item != 2 {
		t.Fatalf("listItem(-2) = (%d, %v), want (2, true)", item, ok)
	}

	_, ok = listItem([]int{1, 2, 3}, 3)
	if ok {
		t.Fatal("expected out-of-bounds to return ok=false")
	}
}

func BenchmarkListAtTypedInt(b *testing.B) {
	io, inChans, outChans := newIO([]string{"data", "idx"}, []string{"res", "err"})
	handler, err := (listAt{}).Create(io, nil)
	if err != nil {
		b.Fatalf("Create returned error: %v", err)
	}

	ctx := b.Context()
	go handler(ctx)

	listMsg := runtime.NewListIntMsg([]int64{1, 2, 3, 4, 5, 6, 7, 8})
	idxMsg := runtime.NewIntMsg(4)

	b.ResetTimer()
	for range b.N {
		inChans["data"] <- runtime.OrderedMsg{Msg: listMsg}
		inChans["idx"] <- runtime.OrderedMsg{Msg: idxMsg}
		<-outChans["res"]
	}
}
