package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// TestFormatIntReceivesInputsConcurrently checks that formatInt does not depend on input arrival order.
func TestFormatIntReceivesInputsConcurrently(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newNamedRuntimeIO([]string{"data", "base"}, []string{"res"})
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

// TestTernaryReceivesInputsConcurrently checks that ternarySelector consumes all inputs concurrently.
func TestTernaryReceivesInputsConcurrently(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newNamedRuntimeIO([]string{"if", "then", "else"}, []string{"res"})
	handler, err := (ternarySelector{}).Create(io, nil)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		handler(ctx)
		close(done)
	}()

	for _, order := range [][]string{{"if", "then", "else"}, {"else", "then", "if"}} {
		sendInOrder(t, inChans, order, map[string]runtime.Msg{
			"if":   runtime.NewBoolMsg(true),
			"then": runtime.NewStringMsg("then"),
			"else": runtime.NewStringMsg("else"),
		})
		assertOutputEquals(t, outChans, "res", runtime.NewStringMsg("then"), order)
	}

	cancel()
	<-done
}
