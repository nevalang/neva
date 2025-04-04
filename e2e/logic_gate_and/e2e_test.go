package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// Check that logical AND works
func TestAND(t *testing.T) {
	cmd := exec.Command("neva", "run", "main")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	require.Equal(
		t,
		"true\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
