package e2e

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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
	stdin         string
	captureStderr bool
	expectedCode  int
}

// WithStdin sets stdin input for the command.
func WithStdin(stdin string) Option {
	return func(c *config) {
		c.stdin = stdin
	}
}

// WithStderr captures both stdout and stderr together (combined output).
func WithStderr() Option {
	return func(c *config) {
		c.captureStderr = true
	}
}

// WithCode sets the expected exit code (default is 0).
func WithCode(code int) Option {
	return func(c *config) {
		c.expectedCode = code
	}
}

// Run executes the neva command with the given arguments and options.
// It returns captured output (stdout only by default, or combined stdout+stderr if WithStderr is set).
// The working directory is os.Getwd() (relies on go test running each package with cwd at that package).
func Run(t *testing.T, args []string, opts ...Option) string {
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

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if cfg.captureStderr {
		// Capture combined output - both stdout and stderr go to the same buffer
		writer := io.MultiWriter(&stdout, &stderr)
		cmd.Stdout = writer
		cmd.Stderr = writer
	} else {
		// Capture stdout only, stderr goes to separate buffer (discarded unless error)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	}

	err = cmd.Run()
	actualCode := getExitCode(err)

	// Always show both stdout and stderr in error messages for consistency
	// When captureStderr is true, both buffers contain the same combined content (via MultiWriter)
	outputMsg := fmt.Sprintf("stdout: %q\nstderr: %q", stdout.String(), stderr.String())
	require.Equal(t, cfg.expectedCode, actualCode,
		"neva execution exit code mismatch. %s", outputMsg)

	// Return captured output:
	// - When captureStderr is true: stdout contains combined output (stderr buffer has same content via MultiWriter)
	// - When captureStderr is false: stdout contains only stdout
	return stdout.String()
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
