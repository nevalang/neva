package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime/messages"
)

// eq_test.go contains unit tests for eq runtime function.

// TestEqProducesExpectedValue checks equality behavior.
func TestEqProducesExpectedValue(t *testing.T) {
	t.Parallel()
	assertBinaryOperatorResult(t, eq{}, messages.NewStringMsg("same"), messages.NewStringMsg("same"), messages.NewBoolMsg(true))
}
