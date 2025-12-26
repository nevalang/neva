package test

import (
	"os"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

// This function tests `neva build` generates a Windows executable.
func TestBuildWindows(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
	}()

	e2e.Run(t, "new", ".")

	out := e2e.RunCombined(t, "build", "--target-os=windows", "--target-arch=amd64", "src")
	defer func() {
		require.NoError(t, os.Remove("output.exe"))
	}()

	_, err := os.Stat("output.exe")
	require.NoError(t, err, out)
}
