package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Test successful case with "Alice"
	out, _ := e2e.Run(t, []string{"run", "switch"}, e2e.WithStdin("Alice\n"))
	require.Equal(t, "Enter the name: ALICE\n", out)

	// Test panic case with "Bob"
	out, _ = e2e.Run(t, []string{"run", "switch"}, e2e.WithStdin("Bob\n"))
	require.Equal(t, "Enter the name: bob\n", out)

	// Test panic case with "Charlie"
	out, stderr := e2e.Run(t, []string{"run", "switch"}, e2e.WithStdin("Charlie\n"))
	require.Equal(t, "Enter the name: panic: Charlie\n", out+stderr)
}
