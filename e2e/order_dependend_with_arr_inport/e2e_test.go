package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// we do 100 attempts because there was a floating bug
	// caused by unordered map that was
	for i := 0; i < 100; i++ {
		cmd := exec.Command("neva", "run", "main")

		out, err := cmd.CombinedOutput()
		require.NoError(t, err)
		require.Equal(
			t,
			"-4\n",
			string(out),
		)

		require.Equal(t, 0, cmd.ProcessState.ExitCode())
	}
}
