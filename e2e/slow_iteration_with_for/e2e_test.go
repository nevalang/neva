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
		"1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
