package e2e

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	nevaos "github.com/nevalang/neva/pkg/os"
	"github.com/stretchr/testify/require"
)

// Option configures Run behavior.
type Option func(*config)

type config struct {
	stdin        string
	wd           string
	expectedCode int
	timeout      time.Duration
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

// WithTimeout sets the maximum command execution time.
func WithTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.timeout = timeout
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

	// Respect explicit per-test override; otherwise derive a safe default.
	runTimeout := resolveRunTimeout(t, cfg.timeout)
	ctx, cancel := context.WithTimeout(context.Background(), runTimeout)
	defer cancel()

	repoRoot := FindRepoRoot(t)
	mainPath := filepath.Join(repoRoot, "cmd", "neva", "main.go")

	// Build the CLI binary from repo root; run it from wd.
	binPath := buildNevaBinary(t, repoRoot, mainPath)

	cmdArgs := append([]string{binPath}, args...)
	// #nosec G204 -- test helper executes commands constructed from test inputs
	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	configureCommandCleanup(cmd)

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
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		require.FailNow(
			t,
			"neva execution timed out",
			"timeout: %s\nargs: %v\nstdout: %q\nstderr: %q",
			runTimeout,
			args,
			stdoutBuf.String(),
			stderrBuf.String(),
		)
	}
	actualCode := getExitCode(err)

	// Always show both stdout and stderr in error messages
	outputMsg := fmt.Sprintf("stdout: %q\nstderr: %q", stdoutBuf.String(), stderrBuf.String())
	require.Equal(t, cfg.expectedCode, actualCode,
		"neva execution exit code mismatch. %s", outputMsg)

	return stdoutBuf.String(), stderrBuf.String()
}

// resolveRunTimeout calculates command timeout for e2e.Run.
// Example: with `go test -timeout=5m`, each command gets at most 60s by default.
func resolveRunTimeout(t *testing.T, configured time.Duration) time.Duration {
	t.Helper()

	// Explicit option wins (`e2e.WithTimeout(...)`).
	if configured > 0 {
		return configured
	}

	const (
		defaultTimeout = 60 * time.Second
		safetyMargin   = 2 * time.Second
		minTimeout     = 5 * time.Second
	)

	deadline, ok := t.Deadline()
	if !ok {
		return defaultTimeout
	}

	// Keep per-run timeouts short even when overall test timeout is large.
	remaining := time.Until(deadline) - safetyMargin
	if remaining > defaultTimeout {
		return defaultTimeout
	}

	// Avoid zero/negative values when the test is close to deadline.
	if remaining < minTimeout {
		return minTimeout
	}

	return remaining
}

// FindRepoRoot finds the repository root using go env GOMOD.
func FindRepoRoot(tb testing.TB) string {
	tb.Helper()

	// #nosec G204 -- command arguments are constant
	cmd := exec.Command("go", "env", "GOMOD")
	output, err := cmd.Output()
	require.NoError(tb, err, "failed to run 'go env GOMOD'")

	gomodPath := strings.TrimSpace(string(output))
	require.NotEmpty(tb, gomodPath, "GOMOD path is empty")

	repoRoot := filepath.Dir(gomodPath)
	require.NotEmpty(tb, repoRoot, "repo root is empty")

	return repoRoot
}

// BuildNevaBinary builds the neva CLI binary from repo root and returns its path.
func BuildNevaBinary(tb testing.TB, repoRoot string) string {
	tb.Helper()

	mainPath := filepath.Join(repoRoot, "cmd", "neva", "main.go")
	return buildNevaBinary(tb, repoRoot, mainPath)
}

// PrepareIsolatedNevaHome creates an isolated Neva home and wires the local stdlib into it.
func PrepareIsolatedNevaHome(repoRoot, homeDir string) error {
	nevaHome := filepath.Join(homeDir, "neva")
	if err := os.MkdirAll(nevaHome, 0o755); err != nil {
		return err
	}

	stdSrc := filepath.Join(repoRoot, "std")
	stdDst := filepath.Join(nevaHome, "std")
	if err := os.Symlink(stdSrc, stdDst); err == nil {
		return nil
	}

	return nevaos.CopyDir(stdSrc, stdDst)
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
func buildNevaBinary(tb testing.TB, repoRoot, mainPath string) string {
	tb.Helper()

	binPath, err := buildNevaBinaryFromCache(repoRoot, mainPath)
	if err == nil {
		return binPath
	}

	tb.Logf("e2e: shared neva binary cache unavailable (%v), falling back to per-test build", err)

	return buildNevaBinaryPerTest(tb, repoRoot, mainPath)
}

// buildNevaBinaryPerTest builds an isolated neva binary for a single test.
func buildNevaBinaryPerTest(tb testing.TB, repoRoot, mainPath string) string {
	tb.Helper()

	binPath := filepath.Join(tb.TempDir(), "neva")
	buildCmd := exec.Command("go", "build", "-o", binPath, mainPath)
	buildCmd.Dir = repoRoot

	var buildStdoutBuf bytes.Buffer
	var buildStderrBuf bytes.Buffer
	buildCmd.Stdout = &buildStdoutBuf
	buildCmd.Stderr = &buildStderrBuf

	err := buildCmd.Run()
	require.NoError(
		tb,
		err,
		"failed to build neva CLI. stdout: %q stderr: %q",
		buildStdoutBuf.String(),
		buildStderrBuf.String(),
	)

	return binPath
}

// buildNevaBinaryFromCache returns a shared neva CLI binary keyed by compiler input fingerprint.
// "Process-safe" means concurrent go test package binaries coordinate with a lock file:
// one process builds, others wait for the artifact (up to lock timeout) and then reuse it.
func buildNevaBinaryFromCache(repoRoot, mainPath string) (string, error) {
	cacheRoot, err := e2eCacheRootDir()
	if err != nil {
		return "", fmt.Errorf("resolve e2e cache root: %w", err)
	}

	cacheKey, err := nevaBuildFingerprint(repoRoot)
	if err != nil {
		return "", fmt.Errorf("compute neva build fingerprint: %w", err)
	}

	buildDir := filepath.Join(cacheRoot, cacheKey)
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return "", fmt.Errorf("create build cache directory: %w", err)
	}

	binPath := filepath.Join(buildDir, nevaBinaryName())
	if nevaos.FileExists(binPath) {
		return binPath, nil
	}

	release, err := acquireCacheLock(filepath.Join(buildDir, ".lock"), 45*time.Second)
	if err != nil {
		return "", fmt.Errorf("acquire build cache lock: %w", err)
	}
	defer release()

	if nevaos.FileExists(binPath) {
		return binPath, nil
	}

	tmpBinPath := filepath.Join(buildDir, fmt.Sprintf("neva.%d.tmp", os.Getpid()))
	if err := buildNevaBinaryToPath(repoRoot, mainPath, tmpBinPath); err != nil {
		return "", err
	}

	if err := os.Rename(tmpBinPath, binPath); err != nil {
		_ = os.Remove(tmpBinPath)
		return "", fmt.Errorf("promote temporary build artifact: %w", err)
	}

	return binPath, nil
}

// buildNevaBinaryToPath compiles cmd/neva into the requested output path.
func buildNevaBinaryToPath(repoRoot, mainPath, binPath string) error {
	buildCmd := exec.Command("go", "build", "-o", binPath, mainPath)
	buildCmd.Dir = repoRoot

	var buildStdoutBuf bytes.Buffer
	var buildStderrBuf bytes.Buffer
	buildCmd.Stdout = &buildStdoutBuf
	buildCmd.Stderr = &buildStderrBuf

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf(
			"failed to build neva CLI. stdout: %q stderr: %q: %w",
			buildStdoutBuf.String(),
			buildStderrBuf.String(),
			err,
		)
	}

	return nil
}

// e2eCacheRootDir returns the directory for cross-process e2e build artifacts.
// It prefers the user cache directory and falls back to the OS temp directory.
func e2eCacheRootDir() (string, error) {
	baseDir, err := os.UserCacheDir()
	if err != nil {
		baseDir = os.TempDir()
	}

	cacheRoot := filepath.Join(baseDir, "neva", "e2e-binary-cache")
	if err := os.MkdirAll(cacheRoot, 0o755); err != nil {
		return "", fmt.Errorf("create cache root directory: %w", err)
	}

	return cacheRoot, nil
}

// nevaBuildFingerprint computes a stable cache key from local build input metadata:
// relative path + size + mtime(ns) for each selected file.
// This path is for on-disk repo files only; stdlib extraction uses content hashing
// because embed.FS metadata does not provide reliable mtimes.
func nevaBuildFingerprint(repoRoot string) (string, error) {
	files, err := compilerInputFiles(repoRoot)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			return "", fmt.Errorf("stat %s: %w", file, err)
		}

		rel, err := filepath.Rel(repoRoot, file)
		if err != nil {
			return "", fmt.Errorf("resolve relative path for %s: %w", file, err)
		}

		normalizedRel := filepath.ToSlash(rel)
		_, _ = fmt.Fprintf(hash, "%s|%d|%d\n", normalizedRel, info.Size(), info.ModTime().UTC().UnixNano())
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// compilerInputFiles returns local files that can affect `go build ./cmd/neva`.
// It includes only repo-local package files from `go list -deps -json` plus go.mod/go.sum.
func compilerInputFiles(repoRoot string) ([]string, error) {
	// Use package metadata to include only files that affect cmd/neva build.
	// External module changes are already covered by go.mod/go.sum.
	cmd := exec.Command("go", "list", "-deps", "-json", "./cmd/neva")
	cmd.Dir = repoRoot

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf(
			"list cmd/neva dependencies: %w (stderr: %s)",
			err,
			strings.TrimSpace(stderr.String()),
		)
	}

	filesSet := map[string]struct{}{}
	decoder := json.NewDecoder(bytes.NewReader(stdout.Bytes()))
	normalizedRoot := filepath.Clean(repoRoot)

	for {
		var payload map[string]json.RawMessage
		if err := decoder.Decode(&payload); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("decode go list output: %w", err)
		}

		dir, shouldInclude, err := decodePackageDir(payload, normalizedRoot)
		if err != nil {
			return nil, err
		}
		if !shouldInclude {
			continue
		}

		for _, field := range []string{"GoFiles", "CgoFiles", "SFiles", "SysoFiles", "EmbedFiles"} {
			if err := addPackageFilesFromRaw(filesSet, dir, payload[field]); err != nil {
				return nil, err
			}
		}
	}

	filesSet[filepath.Join(repoRoot, "go.mod")] = struct{}{}
	goSum := filepath.Join(repoRoot, "go.sum")
	if nevaos.FileExists(goSum) {
		filesSet[goSum] = struct{}{}
	}

	files := make([]string, 0, len(filesSet))
	for file := range filesSet {
		files = append(files, file)
	}
	sort.Strings(files)

	return files, nil
}

// addPackageFiles merges package-relative file names into an absolute file set.
func addPackageFiles(filesSet map[string]struct{}, dir string, names []string) {
	for _, name := range names {
		if name == "" {
			continue
		}
		filesSet[filepath.Join(dir, name)] = struct{}{}
	}
}

// decodePackageDir extracts package directory info and filters out stdlib/external packages.
func decodePackageDir(payload map[string]json.RawMessage, normalizedRoot string) (dir string, include bool, err error) {
	standard := false
	if raw, ok := payload["Standard"]; ok && len(raw) > 0 {
		if err := json.Unmarshal(raw, &standard); err != nil {
			return "", false, fmt.Errorf("decode package standard flag: %w", err)
		}
	}
	if standard {
		return "", false, nil
	}

	rawDir, ok := payload["Dir"]
	if !ok || len(rawDir) == 0 {
		return "", false, nil
	}
	if err := json.Unmarshal(rawDir, &dir); err != nil {
		return "", false, fmt.Errorf("decode package directory: %w", err)
	}
	if dir == "" {
		return "", false, nil
	}

	normalizedDir := filepath.Clean(dir)
	if normalizedDir != normalizedRoot &&
		!strings.HasPrefix(normalizedDir, normalizedRoot+string(os.PathSeparator)) {
		return "", false, nil
	}

	return dir, true, nil
}

// addPackageFilesFromRaw decodes `go list -json` file arrays and merges them into the file set.
func addPackageFilesFromRaw(filesSet map[string]struct{}, dir string, raw json.RawMessage) error {
	if len(raw) == 0 {
		return nil
	}

	var names []string
	if err := json.Unmarshal(raw, &names); err != nil {
		return fmt.Errorf("decode go list file set for %s: %w", dir, err)
	}

	addPackageFiles(filesSet, dir, names)

	return nil
}

// acquireCacheLock acquires an exclusive lock file and returns a release function.
// The lock is polled until timeout; the returned function must be called to close and remove the lock file.
func acquireCacheLock(lockPath string, timeout time.Duration) (func(), error) {
	deadline := time.Now().Add(timeout)

	// Repeatedly try to become the lock owner until timeout expires.
	for {
		lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0o600)
		if err == nil {
			return func() {
				_ = lockFile.Close()
				_ = os.Remove(lockPath)
			}, nil
		}

		if !errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("create lock file: %w", err)
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timeout waiting for lock: %s", lockPath)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// nevaBinaryName returns platform-correct executable name for the cached CLI binary.
func nevaBinaryName() string {
	if os.PathSeparator == '\\' {
		return "neva.exe"
	}

	return "neva"
}
