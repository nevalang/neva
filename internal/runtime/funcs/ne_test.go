package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// ne_test.go contains unit tests for notEq runtime function.

// TestNotEqProducesExpectedValue checks inequality behavior.
func TestNotEqProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertBinaryOperatorResult(t, notEq{}, runtime.NewIntMsg(1), runtime.NewIntMsg(2), runtime.NewBoolMsg(true))
}
