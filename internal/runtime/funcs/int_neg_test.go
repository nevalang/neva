package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime/messages"
)

// int_neg_test.go contains unit tests for intNeg runtime function.

// TestIntNegProducesExpectedValue checks negation behavior.
func TestIntNegProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertUnaryOperatorResult(t, intNeg{}, messages.NewIntMsg(8), messages.NewIntMsg(-8))
}
