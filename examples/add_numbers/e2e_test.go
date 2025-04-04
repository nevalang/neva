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
		cmd := exec.Command("neva", "run", "add_numbers")

		out, err := cmd.CombinedOutput()
		require.NoError(t, err, string(out))
		require.Equal(
			t,
			"3\n",
			string(out),
		)

		require.Equal(t, 0, cmd.ProcessState.ExitCode())
	}
}
