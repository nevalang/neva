package test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	for i := 0; i < 1; i++ {
		t.Run("Add numbers from stdin", func(t *testing.T) {
			cmd := exec.Command("neva", "run", "main")

			cmd.Stdin = strings.NewReader("3\n4\n\n")

			out, err := cmd.CombinedOutput()
			require.NoError(t, err)

			require.Equal(
				t,
				"7\n",
				string(out),
			)

			require.Equal(t, 0, cmd.ProcessState.ExitCode())
		})
	}
}
