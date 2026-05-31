package funcs

import (
	"testing"
	"time"

	"github.com/nevalang/neva/internal/runtime"
)

// TestListAtSupportsNegativeIndex verifies invariant: -1 returns the last element.
func TestListAtSupportsNegativeIndex(t *testing.T) {
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

	inChans["data"] <- runtime.OrderedMsg{Msg: runtime.NewListMsg([]runtime.Msg{runtime.NewIntMsg(10), runtime.NewIntMsg(20), runtime.NewIntMsg(30)})}
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

// TestListAtOutOfBounds verifies invariant: out-of-bounds index emits `err`.
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

	inChans["data"] <- runtime.OrderedMsg{Msg: runtime.NewListMsg([]runtime.Msg{runtime.NewStringMsg("x")})}
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
