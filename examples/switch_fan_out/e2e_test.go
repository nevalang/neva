package test

import (
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// Test successful case with "Alice"
	out := e2e.RunExampleWithStdin(t, "Alice\n", "switch_fan_out")
	require.Equal(t, "Enter the name: ALICEalice\n", out)

	// Test panic case with "Bob"
	// RunWithStdinCombined captures stdout+stderr so the panic text emitted by
	// runtime.Panic remains part of the returned output.
	out = e2e.RunExampleWithStdinCombined(t, "Bob\n", "switch_fan_out")
	require.Equal(t, "Enter the name: panic: Bob\n", out)
}
