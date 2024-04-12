package test

import (
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"testing"
)

func Test(t *testing.T) {
	err := os.Chdir("../")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "18_split_join_string")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"neva\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
