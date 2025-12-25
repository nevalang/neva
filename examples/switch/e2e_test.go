package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Test successful case with "Alice"
	out := e2e.RunExampleWithStdin(t, "Alice\n", "switch")
	require.Equal(t, "Enter the name: ALICE\n", out)

	// Test panic case with "Bob"
	out = e2e.RunExampleWithStdin(t, "Bob\n", "switch")
	require.Equal(t, "Enter the name: bob\n", out)

	// Test panic case with "Charlie"
	out = e2e.RunExampleWithStdinCombined(t, "Charlie\n", "switch")
	require.Equal(t, "Enter the name: panic: Charlie\n", out)
}
