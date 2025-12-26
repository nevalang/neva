package e2e

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Option configures Run behavior.
type Option func(*config)

type config struct {
	stdin        string
	expectedCode int
}

// WithStdin sets stdin input for the command.
func WithStdin(stdin string) Option {
	return func(c *config) {
		c.stdin = stdin
	}
}

// WithCode sets the expected exit code (default is 0).
func WithCode(code int) Option {
	return func(c *config) {
		c.expectedCode = code
	}
}

// Run executes the neva command with the given arguments and options.
// It returns captured stdout and stderr separately.
// The working directory is os.Getwd() (relies on go test running each package with cwd at that package).
// Most tests can ignore stderr: `out, _ := e2e.Run(...)`
// Tests that need stderr (e.g., panic cases) can use: `out, stderr := e2e.Run(...)` and combine if needed.
func Run(t *testing.T, args []string, opts ...Option) (stdout, stderr string) {
	t.Helper()

	cfg := &config{
		expectedCode: 0,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	repoRoot := findRepoRoot(t)
	mainPath := filepath.Join(repoRoot, "cmd", "neva", "main.go")
	wd, err := os.Getwd()
	require.NoError(t, err, "failed to get working directory")

	// For examples, we need to run from the examples directory to find neva.yml
	// Check if we're in an examples subdirectory and adjust working directory accordingly
	examplesDir := filepath.Join(repoRoot, "examples")
	if strings.HasPrefix(wd, examplesDir+string(filepath.Separator)) || wd == examplesDir {
		// We're in an examples subdirectory, run from examples directory
		wd = examplesDir
	}

	cmdArgs := append([]string{"run", mainPath}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = wd

	if cfg.stdin != "" {
		cmd.Stdin = strings.NewReader(cfg.stdin)
	}

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	// Always capture stdout and stderr separately
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()
	actualCode := getExitCode(err)

	// Always show both stdout and stderr in error messages
	outputMsg := fmt.Sprintf("stdout: %q\nstderr: %q", stdoutBuf.String(), stderrBuf.String())
	require.Equal(t, cfg.expectedCode, actualCode,
		"neva execution exit code mismatch. %s", outputMsg)

	return stdoutBuf.String(), stderrBuf.String()
}

// findRepoRoot finds the repository root using go env GOMOD.
func findRepoRoot(t *testing.T) string {
	t.Helper()

	cmd := exec.Command("go", "env", "GOMOD")
	output, err := cmd.Output()
	require.NoError(t, err, "failed to run 'go env GOMOD'")

	gomodPath := strings.TrimSpace(string(output))
	require.NotEmpty(t, gomodPath, "GOMOD path is empty")

	repoRoot := filepath.Dir(gomodPath)
	require.NotEmpty(t, repoRoot, "repo root is empty")

	return repoRoot
}

// getExitCode extracts the exit code from an error.
// Returns 0 if err is nil, otherwise extracts from *exec.ExitError.
func getExitCode(err error) int {
	if err == nil {
		return 0
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return exitErr.ExitCode()
	}

	// If it's not an ExitError, we can't determine the exit code
	// This shouldn't happen in normal operation, but return a non-zero
	// to indicate failure
	return 1
}
