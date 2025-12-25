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

// RunInDir is the same as Run but sets the working directory for the invoked command.
func RunInDir(t *testing.T, dir string, args ...string) string {
	t.Helper()
	return runWithModeInDir(t, dir, "", captureStdoutOnly, args...)
}

// RunWithStdin executes the neva command with the given arguments and stdin input.
// It asserts success (exit code 0) and returns captured stdout.
// Stderr is passed through to os.Stderr for visibility on failure.
func RunWithStdin(t *testing.T, stdin string, args ...string) string {
	t.Helper()
	return runWithMode(t, stdin, captureStdoutOnly, args...)
}

// RunWithStdinInDir is the same as RunWithStdin but sets the working directory for the invoked command.
func RunWithStdinInDir(t *testing.T, dir, stdin string, args ...string) string {
	t.Helper()
	return runWithModeInDir(t, dir, stdin, captureStdoutOnly, args...)
}

// RunCombined executes the neva command and captures both stdout and stderr (combined),
// asserting a zero exit code. Combined output is returned as a string.
func RunCombined(t *testing.T, args ...string) string {
	t.Helper()
	return runWithMode(t, "", captureCombinedOutput, args...)
}

// RunCombinedInDir is the same as RunCombined but sets the working directory for the invoked command.
func RunCombinedInDir(t *testing.T, dir string, args ...string) string {
	t.Helper()
	return runWithModeInDir(t, dir, "", captureCombinedOutput, args...)
}

// RunWithStdinCombined is similar to RunCombined but also lets callers pass stdin.
func RunWithStdinCombined(t *testing.T, stdin string, args ...string) string {
	t.Helper()
	return runWithMode(t, stdin, captureCombinedOutput, args...)
}

// RunWithStdinCombinedInDir is similar to RunWithStdinCombined but also sets the working directory.
func RunWithStdinCombinedInDir(t *testing.T, dir, stdin string, args ...string) string {
	t.Helper()
	return runWithModeInDir(t, dir, stdin, captureCombinedOutput, args...)
}

// RepoRoot returns the repository root directory (the directory containing cmd/neva/main.go).
func RepoRoot(t *testing.T) string {
	t.Helper()
	return repoRoot(t)
}

// ExamplesDir returns the absolute path to the repository examples directory.
func ExamplesDir(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "examples")
}

// RunExample runs `neva run <exampleName>` from the repository `examples/` directory.
func RunExample(t *testing.T, exampleName string) string {
	t.Helper()
	return RunInDir(t, ExamplesDir(t), "run", exampleName)
}

// RunExampleCombined runs `neva run <exampleName>` from the repository `examples/` directory,
// capturing both stdout and stderr.
func RunExampleCombined(t *testing.T, exampleName string) string {
	t.Helper()
	return RunCombinedInDir(t, ExamplesDir(t), "run", exampleName)
}

// RunExampleWithStdin runs `neva run <exampleName>` from the repository `examples/` directory with stdin.
func RunExampleWithStdin(t *testing.T, stdin, exampleName string) string {
	t.Helper()
	return RunWithStdinInDir(t, ExamplesDir(t), stdin, "run", exampleName)
}

// RunExampleWithStdinCombined runs `neva run <exampleName>` from the repository `examples/` directory with stdin,
// capturing both stdout and stderr.
func RunExampleWithStdinCombined(t *testing.T, stdin, exampleName string) string {
	t.Helper()
	return RunWithStdinCombinedInDir(t, ExamplesDir(t), stdin, "run", exampleName)
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
	return runWithModeInDir(t, "", stdin, mode, args...)
}

func runWithModeInDir(t *testing.T, dir, stdin string, mode outputMode, args ...string) string {
	t.Helper()

	cmd := buildGoRunCmd(t, args...)
	if dir != "" {
		cmd.Dir = dir
	}
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
	root := repoRoot(t)
	main := filepath.Join(root, "cmd", "neva", "main.go")
	cmdArgs := append([]string{"run", main}, args...)

	return exec.Command("go", cmdArgs...)
}

func repoRoot(t *testing.T) string {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)

	return filepath.Join(filepath.Dir(filename), "..", "..")
}
