package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	// we do 10 iterations because there was a bug
	// that was only reproducible after a few runs
	for i := 0; i < 10; i++ {
		cmd := exec.Command("neva", "run", "main")

		out, err := cmd.CombinedOutput()
		require.NoError(t, err)

		require.Equal(
			t,
			"[49,29,19,99]\n",
			string(out),
			"iteration %d", i,
		)

		require.Equal(t, 0, cmd.ProcessState.ExitCode())
	}
}
