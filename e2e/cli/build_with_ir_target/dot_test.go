package test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildIRDOT(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("program.dot"))
	}()

	// create new project
	cmd := exec.Command("neva", "new", ".")
	require.NoError(t, cmd.Run())

	// build with ir target and dot format
	cmd = exec.Command("neva", "build", "-target=ir", "-target-ir-format=dot", "src")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command failed: %v", string(out))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// verify dot file exists
	dotBytes, err := os.ReadFile("program.dot")
	require.NoError(t, err, string(out))

	// Check for basic DOT syntax and content
	dotContent := string(dotBytes)
	require.True(t, strings.HasPrefix(dotContent, "digraph G {"), "DOT file should start with 'digraph G {'")
	require.True(t, strings.Contains(dotContent, "}"), "DOT file should contain closing brace")
}

