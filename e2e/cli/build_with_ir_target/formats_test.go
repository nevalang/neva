package test

import (
	"os"
	"strings"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func TestBuildIRMermaid(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
		require.NoError(t, os.Remove("ir.md"))
	}()

	// create new project
	e2e.Run(t, []string{"new", "."})

	// build with ir target and mermaid format
	out, _ := e2e.Run(t, []string{"build", "-target=ir", "-target-ir-format=mermaid", "src"})

	// verify mermaid file exists
	mdBytes, err := os.ReadFile("ir.md")
	require.NoError(t, err, out)

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
	e2e.Run(t, []string{"new", "."})

	// build with ir target and threejs format
	out, _ := e2e.Run(t, []string{"build", "-target=ir", "-target-ir-format=threejs", "src"})

	// verify threejs file exists
	htmlBytes, err := os.ReadFile("ir.threejs.html")
	require.NoError(t, err, out)

	// Check for basic HTML/ThreeJS content
	htmlContent := string(htmlBytes)
	require.True(t, strings.Contains(htmlContent, "<!DOCTYPE html>"), "File should be HTML")
	require.True(t, strings.Contains(htmlContent, "import * as THREE"), "File should import Three.js")
}

