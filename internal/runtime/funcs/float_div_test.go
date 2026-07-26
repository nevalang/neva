package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime/messages"
)

// float_div_test.go contains unit tests for floatDiv runtime function.

// TestFloatDivProducesExpectedValue checks arithmetic behavior.
func TestFloatDivProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertBinaryOperatorResult(t, floatDiv{}, messages.NewFloatMsg(9.0), messages.NewFloatMsg(2.0), messages.NewFloatMsg(4.5))
}
