package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Test successful case with "Alice"
	out := e2e.Run(t, []string{"run", "switch_fan_out"}, e2e.WithStdin("Alice\n"))
	require.Equal(t, "Enter the name: ALICEalice\n", out)

	// Test panic case with "Bob"
	// RunWithStdinCombined captures stdout+stderr so the panic text emitted by
	// runtime.Panic remains part of the returned output.
	out = e2e.Run(t, []string{"run", "switch_fan_out"}, e2e.WithStdin("Bob\n"), e2e.WithStderr())
	require.Equal(t, "Enter the name: panic: Bob\n", out)
}
