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

	cmd := exec.Command("neva", "run", "stream_to_list")

	// TODO betterh check in a loop
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(
		t,
		"[1,2,3,4,5,6,7,8,9,10]\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
