package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// int_dec_test.go contains unit tests for intDec runtime function.

// TestIntDecProducesExpectedValue checks decrement behavior.
func TestIntDecProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertUnaryOperatorResult(t, intDec{}, runtime.NewIntMsg(41), runtime.NewIntMsg(40))
}
