package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("../..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "2_hello_world/2_with_then_connection")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(
		t,
		"Hello, World!\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
