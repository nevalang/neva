package test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "image/minimal_png")

	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// Check file exists.
	const filename = "minimal.png"

	_, err = os.ReadFile(filename)
	require.NoError(t, err)

	// Remove file output.
	os.Remove(filename)
}
