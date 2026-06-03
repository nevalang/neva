package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// format_int_test.go contains unit tests for formatInt runtime function.

// TestFormatIntReceivesInputsConcurrently verifies order-independent input handling.
func TestFormatIntReceivesInputsConcurrently(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newIO([]string{"data", "base"}, []string{"res"})
	handler, err := (formatInt{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()

	for _, order := range [][]string{{"data", "base"}, {"base", "data"}} {
		sendInOrder(t, inChans, order, map[string]runtime.Msg{
			"data": runtime.NewIntMsg(42),
			"base": runtime.NewIntMsg(10),
		})
		assertOutputEquals(t, outChans, "res", runtime.NewStringMsg("42"), order)
	}

	cancel()
	<-done
}
