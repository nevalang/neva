package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Run executes the neva command with the given arguments.
// It asserts success (exit code 0) and returns captured stdout.
// Stderr is passed through to os.Stderr for visibility on failure.
func Run(t *testing.T, args ...string) string {
	t.Helper()
	return RunWithStdin(t, "", args...)
}

// RunWithStdin executes the neva command with the given arguments and stdin input.
// It asserts success (exit code 0) and returns captured stdout.
// Stderr is passed through to os.Stderr for visibility on failure.
func RunWithStdin(t *testing.T, stdin string, args ...string) string {
	t.Helper()

	cmd := exec.Command("neva", args...)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	require.NoError(t, err, "neva execution failed. Stdout: %s", stdout.String())
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	return stdout.String()
}

// RunExpectingError executes the neva command and asserts it fails with the expected exit code (default 1).
// It returns stdout and stderr captured as strings.
func RunExpectingError(t *testing.T, args ...string) (string, string) {
	t.Helper()

	cmd := exec.Command("neva", args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	require.Error(t, err, "neva execution succeeded but should have failed")
	require.NotEqual(t, 0, cmd.ProcessState.ExitCode(), "exit code should not be 0")

	return stdout.String(), stderr.String()
}
