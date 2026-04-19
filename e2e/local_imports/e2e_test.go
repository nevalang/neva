package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

// Check that math example with multiplying numbers by using port bridge works as expected.
func TestMathMultiplyNumbers(t *testing.T) {
	out, _ := e2e.Run(t, []string{"run", "main"})
	require.Equal(
		t,
		"42\n42\n",
		out,
	)
}
