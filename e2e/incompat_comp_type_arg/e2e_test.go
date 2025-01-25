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
		"main/main.neva:2:1: Subtype must be either union or literal: want int | float, got any\n",
	)
}
