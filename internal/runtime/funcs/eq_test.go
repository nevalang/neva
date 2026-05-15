package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

// eq_test.go contains unit tests for eq runtime function.

// TestEqProducesExpectedValue checks equality behavior.
func TestEqProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertBinaryOperatorResult(t, eq{}, runtime.NewStringMsg("same"), runtime.NewStringMsg("same"), runtime.NewBoolMsg(true))
}
