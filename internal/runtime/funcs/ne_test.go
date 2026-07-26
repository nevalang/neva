package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime/messages"
)

// ne_test.go contains unit tests for notEq runtime function.

// TestNotEqProducesExpectedValue checks inequality behavior.
func TestNotEqProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertBinaryOperatorResult(t, notEq{}, messages.NewIntMsg(1), messages.NewIntMsg(2), messages.NewBoolMsg(true))
}
