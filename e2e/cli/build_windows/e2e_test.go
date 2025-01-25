package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// This function tests `neva build` generates a Windows executable.
func TestBuildWindows(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
	}()

	cmd := exec.Command("neva", "new")
	require.NoError(t, cmd.Run())

	cmd = exec.Command("neva", "build", "--target-os=windows", "--target-arch=amd64", "src")
	require.NoError(t, cmd.Run())
	require.Equal(t, 0, cmd.ProcessState.ExitCode())
	defer func() {
		require.NoError(t, os.Remove("output.exe"))
	}()

	_, err := os.Stat("output.exe")
	require.NoError(t, err)
}
