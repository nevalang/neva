package test

import (
	"strings"
	"testing"
	"time"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	start := time.Now()
	out, _ := e2e.Run(t, []string{"run", "."})
	elapsed := time.Since(start)

	// Check execution time is between 1-10 seconds
	require.GreaterOrEqual(t, elapsed.Seconds(), 1.0)
	require.LessOrEqual(t, elapsed.Seconds(), 10.0)

	// Split output into lines and verify contents
	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.Equal(t, 7, len(lines), out) // Hello + World + 5 numbers

	// First line must be Hello
	require.Equal(t, "Hello", lines[0], out)

	// Check remaining lines contain all expected values
	expected := map[string]bool{
		"World": false,
		"1":     false,
		"2":     false,
		"3":     false,
		"4":     false,
		"5":     false,
	}

	// Check remaining lines contain all expected values
	for _, line := range lines[1:] {
		_, exists := expected[line]
		require.True(t, exists, "Unexpected value in output: %s", line)
		expected[line] = true
	}

	// Verify all expected values were found
	for val, found := range expected {
		require.True(t, found, "Expected value not found: %s", val)
	}
}
