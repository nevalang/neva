package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime/messages"
)

// int_dec_test.go contains unit tests for intDec runtime function.

// TestIntDecProducesExpectedValue checks decrement behavior.
func TestIntDecProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertUnaryOperatorResult(t, intDec{}, messages.NewIntMsg(41), messages.NewIntMsg(40))
}
