package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	cmd := exec.Command("neva", "run", "main")
	cmd.Env = append(os.Environ(), "GOTOOLCHAIN=go1.25.0")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	require.Equal(t, "hello neva lang\n", string(out))
	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
