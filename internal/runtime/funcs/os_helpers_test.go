package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// os_helpers_test.go contains unit tests for shared std/os runtime helpers.

func TestCreateBinaryLoopReceivesInputsConcurrently(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newIO([]string{"first", "second"}, []string{"res", "err"})
	handler, err := createBinaryLoop(io, "first", "second", func(first, second runtime.OrderedMsg) (runtime.Msg, error) {
		return runtime.NewStringMsg(first.Str() + ":" + second.Str()), nil
	})
	if err != nil {
		t.Fatalf("createBinaryLoop returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()

	// Intent: catch sequential two-input receives by sending the second input first.
	sendInOrder(t, inChans, []string{"second", "first"}, map[string]runtime.Msg{
		"first":  runtime.NewStringMsg("left"),
		"second": runtime.NewStringMsg("right"),
	})
	assertOutputEquals(t, outChans, "res", runtime.NewStringMsg("left:right"), []string{"second", "first"})

	cancel()
	<-done
}
