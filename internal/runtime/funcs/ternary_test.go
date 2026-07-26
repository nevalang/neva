package funcs

import (
	"context"
	"testing"

	"github.com/nevalang/neva/internal/runtime/messages"
)

// ternary_test.go contains unit tests for ternarySelector runtime function.

// TestTernaryReceivesInputsConcurrently verifies order-independent input handling.
func TestTernaryReceivesInputsConcurrently(t *testing.T) {
	t.Parallel()

	io, inChans, outChans := newIO([]string{"if", "then", "else"}, []string{"res"})
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
		sendInOrder(t, inChans, order, map[string]messages.Msg{
			"if":   messages.NewBoolMsg(true),
			"then": messages.NewStringMsg("then"),
			"else": messages.NewStringMsg("else"),
		})
		assertOutputEquals(t, outChans, "res", messages.NewStringMsg("then"), order)
	}

	cancel()
	<-done
}
