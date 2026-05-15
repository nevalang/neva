package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// float_div_test.go contains unit tests for floatDiv runtime function.

// TestFloatDivProducesExpectedValue checks arithmetic behavior.
func TestFloatDivProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertBinaryOperatorResult(t, floatDiv{}, runtime.NewFloatMsg(9.0), runtime.NewFloatMsg(2.0), runtime.NewFloatMsg(4.5))
}
