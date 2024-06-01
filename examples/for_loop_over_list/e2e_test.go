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

	// We do 100 attempts to prove that implementation is correct
	// and order of elements is always the same.
	for i := 0; i < 100; i++ {
		cmd := exec.Command("neva", "run", "for_loop_over_list")
		out, err := cmd.CombinedOutput()
		require.NoError(t, err)
		require.Equal(
			t,
			"1\n2\n3\n",
			string(out),
		)
		require.Equal(t, 0, cmd.ProcessState.ExitCode())
	}
}
