package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("../")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "9_const_refs")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"map[l:[1 2 3] m:map[one:1 three:3 two:2]]\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}