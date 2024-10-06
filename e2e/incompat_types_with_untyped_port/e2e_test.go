package test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	cmd := exec.Command("neva", "run", "main")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.True(
		t,
		strings.Contains(
			string(out),
			"Incompatible types: in:data -> println: Subtype inst must have same ref as supertype: got any, want int",
		),
		"Error message should end with expected suffix",
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
