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

	for i := 0; i < 1; i++ {
		cmd := exec.Command("neva", "run", "compare_values")
		out, err := cmd.CombinedOutput()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				t.Fatalf("Program panicked with output:\n%s\nError: %v", string(out), exitError)
			} else {
				t.Fatalf("Error running command: %v", err)
			}
		}

		require.Equal(
			t,
			"They match\n",
			string(out),
		)

		require.Equal(t, 0, cmd.ProcessState.ExitCode())
	}
}
