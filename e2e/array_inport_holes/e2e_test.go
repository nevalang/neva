package test

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	cmd := exec.Command("neva", "run", "main")
	cmd.Stdin = strings.NewReader("yo\n")
	out, _ := cmd.CombinedOutput()
	require.Equal(t, 1, cmd.ProcessState.ExitCode())
	require.Contains(
		t,
		string(out),
		"main/main.neva:4:1: array inport 'printf:args' is used incorrectly: slot 1 is missing\n",
	)
}
