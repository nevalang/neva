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

	for i := 0; i < 100; i++ {
		cmd := exec.Command("neva", "run", "select")
		out, err := cmd.CombinedOutput()
		require.NoError(t, err)
		require.Equal(
			t,
			"a\nb\nc\nd\n",
			string(out),
			"iteration %d failed\n", i,
		)
		require.Equal(t, 0, cmd.ProcessState.ExitCode())
	}
}
