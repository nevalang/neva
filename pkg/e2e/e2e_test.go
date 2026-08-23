package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanGitHead(t *testing.T) {
	repoRoot := t.TempDir()
	testGit(t, repoRoot, "init")
	testGit(t, repoRoot, "config", "user.email", "test@example.com")
	testGit(t, repoRoot, "config", "user.name", "Test User")

	filePath := filepath.Join(repoRoot, "main.go")
	require.NoError(t, os.WriteFile(filePath, []byte("package main\n"), 0o600))
	testGit(t, repoRoot, "add", "main.go")
	testGit(t, repoRoot, "commit", "-m", "initial")

	want := testGit(t, repoRoot, "rev-parse", "HEAD")
	got, ok := cleanGitHead(repoRoot)
	require.True(t, ok)
	require.Equal(t, want, got)

	require.NoError(t, os.WriteFile(filePath, []byte("package main\n// changed\n"), 0o600))
	_, ok = cleanGitHead(repoRoot)
	require.False(t, ok)
}

func testGit(t *testing.T, repoRoot string, args ...string) string {
	t.Helper()

	// #nosec G204 -- test arguments are literals defined in TestCleanGitHead.
	//nolint:noctx // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	cmd := exec.Command("git", args...)
	cmd.Dir = repoRoot

	output, err := cmd.Output()
	require.NoError(t, err)

	return strings.TrimSpace(string(output))
}
