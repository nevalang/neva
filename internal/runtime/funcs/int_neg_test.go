package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// int_neg_test.go contains unit tests for intNeg runtime function.

// TestIntNegProducesExpectedValue checks negation behavior.
func TestIntNegProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertUnaryOperatorResult(t, intNeg{}, runtime.NewIntMsg(8), runtime.NewIntMsg(-8))
}
