package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// Check that math example with multiplying numbers by using port bridge works as expected.
func TestMathMultiplyNumbers(t *testing.T) {
	cmd := exec.Command("neva", "run", "main")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)
	require.Equal(
		t,
		"42\n42\n",
		string(out),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
