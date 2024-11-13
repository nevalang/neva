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
	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Contains(
		t,
		string(out),
		"main/main.neva: array inport 'printf:args' is used incorrectly: slot 1 is missing\n",
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
