package test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("../")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	cmd := exec.Command("go", "run", "../cmd/neva", "run", "20_file_writeall")

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	require.Equal(
		t,
		"",
		strings.TrimSuffix(string(out), "\n"),
	)

	require.Equal(t, 0, cmd.ProcessState.ExitCode())

	// Check file contents.
	const filename = "file_writer_example.txt"

	want, err := os.ReadFile(filename)
	require.NoError(t, err)
	require.Equal(
		t,
		"Hello, io.WriteAll!",
		string(want),
	)

	// Remove file output.
	os.Remove(filename)
}
