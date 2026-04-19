package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Test successful case with "Alice"
	out, _ := e2e.Run(t, []string{"run", "."}, e2e.WithStdin("Alice\n"))
	require.Equal(t, "Enter the name: ALICEalice\n", out)

	// Test panic case with "Bob"
	// Combine stdout+stderr so the panic text emitted by runtime.Panic remains part of the output.
	out, stderr := e2e.Run(t, []string{"run", "."}, e2e.WithStdin("Bob\n"))
	require.Equal(t, "Enter the name: panic: Bob\n", out+stderr)
}
