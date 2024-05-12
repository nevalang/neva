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

	for i := 0; i < 10; i++ {
		cmd := exec.Command("neva", "run", "filter_list")

		out, err := cmd.CombinedOutput()
		require.NoError(t, err)

		require.Equal(
			t,
			"[2,4,6,8,10]\n",
			string(out),
		)

		require.Equal(t, 0, cmd.ProcessState.ExitCode())
	}
}
