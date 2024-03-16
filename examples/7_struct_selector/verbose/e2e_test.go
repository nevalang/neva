package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStructSelectorWithVerbose(t *testing.T) {
	err := os.Chdir("../..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "7_struct_selector/verbose")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(
		t,
		"Charley\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
