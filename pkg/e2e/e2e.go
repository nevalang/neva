package e2e

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"path/filepath"
	"runtime"

	"github.com/stretchr/testify/require"
)

type outputMode int

const (
	captureStdoutOnly outputMode = iota
	captureCombinedOutput
)

// Run executes the neva command with the given arguments.
// It asserts success (exit code 0) and returns captured stdout.
// Stderr is passed through to os.Stderr for visibility on failure.
func Run(t *testing.T, args ...string) string {
	t.Helper()
	return runWithMode(t, "", captureStdoutOnly, args...)
}

// RunWithStdin executes the neva command with the given arguments and stdin input.
// It asserts success (exit code 0) and returns captured stdout.
// Stderr is passed through to os.Stderr for visibility on failure.
func RunWithStdin(t *testing.T, stdin string, args ...string) string {
	t.Helper()
	return runWithMode(t, stdin, captureStdoutOnly, args...)
}

// RunCombined executes the neva command and captures both stdout and stderr (combined),
// asserting a zero exit code. Combined output is returned as a string.
func RunCombined(t *testing.T, args ...string) string {
	t.Helper()
	return runWithMode(t, "", captureCombinedOutput, args...)
}

// RunWithStdinCombined is similar to RunCombined but also lets callers pass stdin.
func RunWithStdinCombined(t *testing.T, stdin string, args ...string) string {
	t.Helper()
	return runWithMode(t, stdin, captureCombinedOutput, args...)
}

// RunExpectingError executes the neva command and asserts it fails with the expected exit code (default 1).
// It returns stdout and stderr captured as strings.
func RunExpectingError(t *testing.T, args ...string) (string, string) {
	t.Helper()

	cmd := buildGoRunCmd(t, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	require.Error(t, err, "neva execution succeeded but should have failed")
	require.NotEqual(t, 0, cmd.ProcessState.ExitCode(), "exit code should not be 0")

	return stdout.String(), stderr.String()
}

func runWithMode(t *testing.T, stdin string, mode outputMode, args ...string) string {
	t.Helper()

	cmd := buildGoRunCmd(t, args...)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	var stdout bytes.Buffer

	switch mode {
	case captureCombinedOutput:
		writer := io.MultiWriter(&stdout, os.Stderr)
		cmd.Stdout = writer
		cmd.Stderr = writer
	default:
		cmd.Stdout = &stdout
		cmd.Stderr = os.Stderr
	}

	err := cmd.Run()
	require.NoError(t, err, "neva execution failed. Stdout/Stderr: %s", stdout.String())
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	return stdout.String()
}

func buildGoRunCmd(t *testing.T, args ...string) *exec.Cmd {
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)

	root := filepath.Join(filepath.Dir(filename), "..", "..")
	main := filepath.Join(root, "cmd", "neva", "main.go")
	cmdArgs := append([]string{"run", main}, args...)

	return exec.Command("go", cmdArgs...)
}
