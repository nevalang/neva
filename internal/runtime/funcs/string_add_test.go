package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime/messages"
)

// string_add_test.go contains unit tests for stringAdd runtime function.

// TestStringAddProducesExpectedValue checks concatenation behavior.
func TestStringAddProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertBinaryOperatorResult(t, stringAdd{}, messages.NewStringMsg("ne"), messages.NewStringMsg("va"), messages.NewStringMsg("neva"))
}
