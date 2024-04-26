package test

// in this file we test files designed specifically for e2e.

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test that CLI will go from up to down and find module's manifest
func Test_UpperThanManifest(t *testing.T) {
	// go one level up (and go back after test is finished)
	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)
	require.NoError(t, os.Chdir(".."))

	cmd := exec.Command("neva", "run", "run_cli_not_from_module_root/foo/bar")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"42\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}

// Test that CLI will go from down to up and find module's manifest
func Test_DownToManifest(t *testing.T) {
	t.Skip() // FIXME https://github.com/nevalang/neva/issues/571

	// go one level down (and go back after test is finished)
	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)
	require.NoError(t, os.Chdir("foo"))

	cmd := exec.Command("neva", "run", "bar")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"42\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
