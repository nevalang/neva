package test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	for i := 0; i < 1; i++ {
		t.Run("Add numbers from stdin", func(t *testing.T) {
			cmd := exec.Command("neva", "run", "add_numbers_from_stdin")

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
