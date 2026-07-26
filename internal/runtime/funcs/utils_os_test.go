package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/messages"
)

// utils_os_test.go contains unit tests for shared std/os runtime helpers.

func TestCreateBinaryLoopReceivesInputsConcurrently(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newIO([]string{"first", "second"}, []string{"res", "err"})
	handler, err := createBinaryLoop(io, "first", "second", func(first, second runtime.OrderedMsg) (messages.Msg, error) {
		return messages.NewStringMsg(first.Str() + ":" + second.Str()), nil
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
	sendInOrder(t, inChans, []string{"second", "first"}, map[string]messages.Msg{
		"first":  messages.NewStringMsg("left"),
		"second": messages.NewStringMsg("right"),
	})
	assertOutputEquals(t, outChans, "res", messages.NewStringMsg("left:right"), []string{"second", "first"})

	cancel()
	<-done
}
