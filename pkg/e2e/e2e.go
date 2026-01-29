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

//nolint:govet // fieldalignment: simple options struct.
type config struct {
	stdin        string
	expectedCode int
	wd           string
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

// WithDir sets the working directory for the command execution.
func WithDir(wd string) Option {
	return func(c *config) {
		c.wd = wd
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

	repoRoot := FindRepoRoot(t)
	mainPath := filepath.Join(repoRoot, "cmd", "neva", "main.go")

	// Build the CLI binary from repo root; run it from wd.
	binPath := buildNevaBinary(t, repoRoot, mainPath)

	cmdArgs := append([]string{binPath}, args...)
	// #nosec G204 -- test helper executes commands constructed from test inputs
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	// Resolve the working directory from which the CLI should be executed.
	// This intentionally differs from the directory used to build the CLI,
	// which must stay at repo root so Go modules resolve correctly.
	wd := cfg.wd
	var err error
	if wd == "" {
		wd, err = os.Getwd()
		require.NoError(t, err, "failed to get working directory")
	}
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

// FindRepoRoot finds the repository root using go env GOMOD.
func FindRepoRoot(t *testing.T) string {
	t.Helper()

	// #nosec G204 -- command arguments are constant
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

// buildNevaBinary builds the neva CLI from the repo root to ensure module
// resolution works regardless of where tests execute the resulting binary.
// It returns the path to the built binary.
func buildNevaBinary(t *testing.T, repoRoot, mainPath string) string {
	t.Helper()

	binPath := filepath.Join(t.TempDir(), "neva")
	buildCmd := exec.Command("go", "build", "-o", binPath, mainPath)
	buildCmd.Dir = repoRoot

	var buildStdoutBuf bytes.Buffer
	var buildStderrBuf bytes.Buffer
	buildCmd.Stdout = &buildStdoutBuf
	buildCmd.Stderr = &buildStderrBuf

	err := buildCmd.Run()
	require.NoError(
		t,
		err,
		"failed to build neva CLI. stdout: %q stderr: %q",
		buildStdoutBuf.String(),
		buildStderrBuf.String(),
	)

	return binPath
}
