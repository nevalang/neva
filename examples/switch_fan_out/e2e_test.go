package test

import (
	"os"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	// Test successful case with "Alice"
	out := e2e.RunWithStdin(t, "Alice\n", "run", "switch_fan_out")
	require.Equal(t, "Enter the name: ALICEalice\n", out)

	// Test panic case with "Bob"
	// RunWithStdinCombined captures stdout+stderr so the panic text emitted by
	// runtime.Panic remains part of the returned output.
	out = e2e.RunWithStdinCombined(t, "Bob\n", "run", "switch_fan_out")
	require.Equal(t, "Enter the name: panic: Bob\n", out)
}
