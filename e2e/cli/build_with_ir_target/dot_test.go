package test

import (
	"os"
	"strings"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func TestBuildIRDOT(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("ir.dot"))
	}()

	// create new project
	e2e.Run(t, "new", ".")

	// build with ir target and dot format
	out := e2e.RunCombined(t, "build", "-target=ir", "-target-ir-format=dot", "src")

	// verify dot file exists
	dotBytes, err := os.ReadFile("ir.dot")
	require.NoError(t, err, out)

	// Check for basic DOT syntax and content
	dotContent := string(dotBytes)
	require.True(t, strings.HasPrefix(dotContent, "digraph G {"), "DOT file should start with 'digraph G {'")
	require.True(t, strings.Contains(dotContent, "}"), "DOT file should contain closing brace")
}

