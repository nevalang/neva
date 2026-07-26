package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime/messages"
)

// int_inc_test.go contains unit tests for intInc runtime function.

// TestIntIncProducesExpectedValue checks increment behavior.
func TestIntIncProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertUnaryOperatorResult(t, intInc{}, messages.NewIntMsg(41), messages.NewIntMsg(42))
}
