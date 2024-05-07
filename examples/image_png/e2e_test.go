package test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	cmd := exec.Command("neva", "run", "image_png")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(
		t,
		"",
		strings.TrimSuffix(string(out), "\n"),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// Check file exists.
	const filename = "minimal.png"

	_, err = os.ReadFile(filename)
	require.NoError(t, err)

	// Remove file output.
	os.Remove(filename)
}
