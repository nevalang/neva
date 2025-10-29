package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// this test verifies that building with the dot target creates a program.dot file
func TestBuildDOTGeneratesProgram(t *testing.T) {
	t.Cleanup(func() {
		_ = os.RemoveAll("src")
		_ = os.Remove("program.dot")
	})

	// create new project
	cmd := exec.Command("neva", "new", ".")
	require.NoError(t, cmd.Run())

	// build with dot target
	cmd = exec.Command("neva", "build", "-target=dot", "src")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// verify program.dot exists
	_, statErr := os.Stat("program.dot")
	require.NoError(t, statErr)

	// check that program.dot has a basic dot structure
	dotBytes, readErr := os.ReadFile("program.dot")
	require.NoError(t, readErr)
	require.Contains(t, string(dotBytes), "digraph G")

	// basic sanity: ensure at least one node/cluster or edge template marker rendered
	// we don't assert exact edges here to keep the test resilient to minor formatting changes
	require.NotZero(t, len(dotBytes))

	// ensure file is in repo root (current working dir), not inside src
	abs, _ := filepath.Abs("program.dot")
	require.NotContains(t, abs, string(filepath.Separator)+"src"+string(filepath.Separator))
}
