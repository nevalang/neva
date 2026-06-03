package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// string_add_test.go contains unit tests for stringAdd runtime function.

// TestStringAddProducesExpectedValue checks concatenation behavior.
func TestStringAddProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertBinaryOperatorResult(t, stringAdd{}, runtime.NewStringMsg("ne"), runtime.NewStringMsg("va"), runtime.NewStringMsg("neva"))
}
