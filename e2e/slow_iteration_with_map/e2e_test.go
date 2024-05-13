package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	cmd := exec.Command("neva", "run", "main")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"[0,1,2,3,4,5,6,7,8,9]\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
