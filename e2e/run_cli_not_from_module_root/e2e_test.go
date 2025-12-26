package test

// in this file we test files designed specifically for e2e.

import (
	"os"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

// Test that CLI will go from up to down and find module's manifest
func Test_UpperThanManifest(t *testing.T) {
	// go one level up (and go back after test is finished)
	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)
	require.NoError(t, os.Chdir(".."))

	out := e2e.Run(t, []string{"run", "run_cli_not_from_module_root/foo/bar"})
	require.Equal(
		t,
		"42\n",
		out,
	)
}

// Test that CLI will go from down to up and find module's manifest
func Test_DownToManifest(t *testing.T) {
	t.Skip() // FIXME https://github.com/nevalang/neva/issues/571

	// go one level down (and go back after test is finished)
	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)
	require.NoError(t, os.Chdir("foo"))

	out := e2e.Run(t, []string{"run", "bar"})
	require.Equal(
		t,
		"42\n",
		out,
	)
}
