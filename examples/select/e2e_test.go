package test

import (
	"fmt"
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
		if err != nil {
			exitError, ok := err.(*exec.ExitError)
			if ok {
				t.Fatalf("Command failed with exit code %d. Error output:\n%s", exitError.ExitCode(), string(out))
			} else {
				t.Fatalf("Command failed with error: %v. Output:\n%s", err, string(out))
			}
		}
		require.Equal(
			t,
			"a\nb\nc\nd\n",
			string(out),
			"iteration %d failed\n", i,
		)
		require.Equal(t, 0, cmd.ProcessState.ExitCode(), fmt.Sprintf("Unexpected exit code on iteration %d", i))
	}
}
