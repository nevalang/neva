package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "wait_group")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	t.Log(string(out))
	require.True(
		t,
		string(out) == "Hello\nWorld!\nNeva\n" ||
			string(out) == "Hello\nNeva\nWorld!\n" ||
			string(out) == "Neva\nHello\nWorld!\n" ||
			string(out) == "Neva\nWorld!\nHello\n" ||
			string(out) == "World!\nHello\nNeva\n" ||
			string(out) == "World!\nNeva\nHello\n",
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
