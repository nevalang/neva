package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

// Check that logical AND works
func TestAND(t *testing.T) {
	out := e2e.Run(t, "run", "main")
	require.Equal(
		t,
		"true\n",
		out,
	)
}
