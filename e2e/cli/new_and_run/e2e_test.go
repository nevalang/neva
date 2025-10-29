package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// This function tests `neva new` followed by `neva run`.
func Test(t *testing.T) {
	defer func() {
		require.NoError(t, os.RemoveAll("src"))
	}()

	cmd := exec.Command("neva", "new", ".")
	require.NoError(t, cmd.Run())

	cmd = exec.Command("neva", "run", "src")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(
		t,
		"Hello, World!\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
