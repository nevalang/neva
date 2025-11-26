package test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildIRMermaid(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("ir.md"))
	}()

	// create new project
	cmd := exec.Command("neva", "new", ".")
	require.NoError(t, cmd.Run())

	// build with ir target and mermaid format
	cmd = exec.Command("neva", "build", "-target=ir", "-target-ir-format=mermaid", "src")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command failed: %v", string(out))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// verify mermaid file exists
	mdBytes, err := os.ReadFile("ir.md")
	require.NoError(t, err, string(out))

	// Check for basic Mermaid syntax and content
	mdContent := string(mdBytes)
	require.True(t, strings.Contains(mdContent, "```mermaid"), "Mermaid file should contain mermaid block")
	require.True(t, strings.Contains(mdContent, "flowchart TD"), "Mermaid file should contain flowchart definition")
}

func TestBuildIRThreeJS(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("ir.threejs.html"))
	}()

	// create new project
	cmd := exec.Command("neva", "new", ".")
	require.NoError(t, cmd.Run())

	// build with ir target and threejs format
	cmd = exec.Command("neva", "build", "-target=ir", "-target-ir-format=threejs", "src")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command failed: %v", string(out))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// verify threejs file exists
	htmlBytes, err := os.ReadFile("ir.threejs.html")
	require.NoError(t, err, string(out))

	// Check for basic HTML/ThreeJS content
	htmlContent := string(htmlBytes)
	require.True(t, strings.Contains(htmlContent, "<!DOCTYPE html>"), "File should be HTML")
	require.True(t, strings.Contains(htmlContent, "import * as THREE"), "File should import Three.js")
}

