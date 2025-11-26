package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

// Check that math example with multiplying numbers by using port bridge works as expected.
func TestMathMultiplyNumbers(t *testing.T) {
	out := e2e.Run(t, "run", "main")
	require.Equal(
		t,
		"6\n",
		out,
	)
}
