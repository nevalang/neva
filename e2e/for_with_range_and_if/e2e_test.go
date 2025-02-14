package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run("", func(t *testing.T) {
			cmd := exec.Command("neva", "run", "main")

			out, err := cmd.CombinedOutput()
			require.NoError(t, err, "out: ", out)

			require.Equal(
				t,
				"1\n0\n",
				string(out),
				"iteration: %d", i,
			)

			require.Equal(t, 0, cmd.ProcessState.ExitCode())
		})
	}
}
