package test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	cmd := exec.Command("neva", "run", "main")
	out, _ := cmd.CombinedOutput()
	require.Equal(t, 1, cmd.ProcessState.ExitCode())
	require.Contains(
		t,
		string(out),
		"Incompatible types: in:data -> println: Subtype inst must have same ref as supertype: got any, want int",
	)
}
