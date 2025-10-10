package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	cmd := exec.Command("neva", "run", "main")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(
		t,
		"50\n30\n20\n100\n",
		string(out),
	)
	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
