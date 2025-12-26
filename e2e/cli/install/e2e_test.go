package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func TestInstall(t *testing.T) {
	workdir := t.TempDir()
	testBinDir := filepath.Join(workdir, "bin")

	// Set GOBIN to our test bin directory
	prevGOBIN := os.Getenv("GOBIN")
	defer func() {
		if prevGOBIN != "" {
			os.Setenv("GOBIN", prevGOBIN)
		} else {
			os.Unsetenv("GOBIN")
		}
	}()

	require.NoError(t, os.Setenv("GOBIN", testBinDir))
	require.NoError(t, os.MkdirAll(testBinDir, 0o755))

	// Create a neva module
	moduleDir := filepath.Join(workdir, "testapp")
	out := e2e.Run(t, []string{"new", moduleDir})
	require.Contains(t, out, "neva module created")

	// Install the program
	srcPath := filepath.Join(moduleDir, "src")
	out = e2e.Run(t, []string{"install", srcPath})

	// Verify the binary was installed
	expectedBinName := "src"
	if runtime.GOOS == "windows" {
		expectedBinName += ".exe"
	}
	installedBinary := filepath.Join(testBinDir, expectedBinName)
	_, err := os.Stat(installedBinary)
	require.NoError(t, err, "binary should exist at %s", installedBinary)

	// Verify the binary is executable
	info, err := os.Stat(installedBinary)
	require.NoError(t, err)
	if runtime.GOOS != "windows" {
		require.NotEqual(t, 0, info.Mode()&0o111, "binary should be executable")
	}

	// Verify the output message contains the expected information
	require.Contains(t, out, "installed")
	require.Contains(t, out, expectedBinName)

	// Verify the command is available from PATH (without full path)
	prevPath := os.Getenv("PATH")
	defer func() {
		os.Setenv("PATH", prevPath)
	}()

	// Add testBinDir to PATH so we can find the command
	newPath := testBinDir + string(filepath.ListSeparator) + prevPath
	require.NoError(t, os.Setenv("PATH", newPath))

	// Verify exec.LookPath can find the command
	foundPath, err := exec.LookPath(expectedBinName)
	require.NoError(t, err, "command should be findable via PATH")
	require.Equal(t, installedBinary, foundPath, "LookPath should return the installed binary path")

	// Run the binary using just the command name (without full path)
	runCmd := exec.Command(expectedBinName)
	runCmd.Dir = workdir
	runOut, runErr := runCmd.CombinedOutput()
	require.NoError(t, runErr, string(runOut))
	require.Contains(t, string(runOut), "Hello, World!")
}

func TestInstallWithPackageDirectory(t *testing.T) {
	workdir := t.TempDir()
	testBinDir := filepath.Join(workdir, "bin")

	// Set GOBIN to our test bin directory
	prevGOBIN := os.Getenv("GOBIN")
	defer func() {
		if prevGOBIN != "" {
			os.Setenv("GOBIN", prevGOBIN)
		} else {
			os.Unsetenv("GOBIN")
		}
	}()

	require.NoError(t, os.Setenv("GOBIN", testBinDir))
	require.NoError(t, os.MkdirAll(testBinDir, 0o755))

	// Create a neva module
	moduleDir := filepath.Join(workdir, "myapp")
	out := e2e.Run(t, []string{"new", moduleDir})
	require.Contains(t, out, "neva module created")

	// Install using the package directory path (not the src subdirectory)
	out = e2e.Run(t, []string{"install", moduleDir})

	// Verify the binary was installed with the correct name (based on package directory)
	expectedBinName := "myapp"
	if runtime.GOOS == "windows" {
		expectedBinName += ".exe"
	}
	installedBinary := filepath.Join(testBinDir, expectedBinName)
	_, err := os.Stat(installedBinary)
	require.NoError(t, err, "binary should exist at %s", installedBinary)

	// Verify the output message
	require.Contains(t, out, "installed")
	require.Contains(t, out, expectedBinName)

	// Verify the command is available from PATH (without full path)
	prevPath := os.Getenv("PATH")
	defer func() {
		os.Setenv("PATH", prevPath)
	}()

	// Add testBinDir to PATH so we can find the command
	newPath := testBinDir + string(filepath.ListSeparator) + prevPath
	require.NoError(t, os.Setenv("PATH", newPath))

	// Verify exec.LookPath can find the command
	foundPath, err := exec.LookPath(expectedBinName)
	require.NoError(t, err, "command should be findable via PATH")
	require.Equal(t, installedBinary, foundPath, "LookPath should return the installed binary path")
}
