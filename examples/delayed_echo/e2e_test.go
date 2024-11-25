package test

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "delayed_echo")

	start := time.Now()
	out, err := cmd.CombinedOutput()
	elapsed := time.Since(start)
	require.NoError(t, err)

	// Check execution time is between 1-5 seconds
	require.GreaterOrEqual(t, elapsed.Seconds(), 1.0)
	require.LessOrEqual(t, elapsed.Seconds(), 5.0)

	// Split output into lines and verify contents
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	require.Equal(t, 7, len(lines)) // Hello + World + 5 numbers

	// First line must be Hello
	require.Equal(t, "Hello", lines[0])

	// Create set of expected remaining values
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

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
