package test

import (
	"strings"
	"testing"
	"time"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// for i := 0; i < 10; i++ {
	start := time.Now()
	out, _ := e2e.Run(t, []string{"run", "."})
	elapsed := time.Since(start)

	// Check execution time is between 1-5 seconds
	require.GreaterOrEqual(t, elapsed.Seconds(), 1.0)
	require.LessOrEqual(t, elapsed.Seconds(), 5.0)

	// Split output into lines and verify contents
	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.Equal(t, 7, len(lines), out) // Hello + World + 5 numbers

	// First line must be Hello
	require.Equal(t, "Hello", lines[0], out)

	// Create set of expected remaining values
	expected := map[string]bool{
		"World": false,
		"1":     false,
		"2":     false,
		"3":     false,
		"4":     false,
		"5":     false,
		"6":     false,
	}

	// Check remaining lines contain all expected values
	// Note: previous test had 6 items in expected map but asserted 7 lines.
	// Output is Hello + World + 5 numbers (1,2,3,4,5).
	// But `delayed_echo` probably outputs 1,2,3,4,5 and World.
	// Previous map had: World, 1, 2, 3, 4, 5. Total 6 items.
	// Hello is separate.
	// Let's copy exact expected map from previous code.
	
	expected = map[string]bool{
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

	// }
}
