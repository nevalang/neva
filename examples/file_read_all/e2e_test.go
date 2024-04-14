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

	cmd := exec.Command("neva", "run", "file_read_all")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	want, err := os.ReadFile("19_file_readall/main.neva")
	require.NoError(t, err)

	require.Equal(
		t,
		string(want),
		strings.TrimSuffix(string(out), "\n"),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())
}
