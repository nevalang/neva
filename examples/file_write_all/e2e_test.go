package test

import (
	"os"
	"strings"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	err := os.Chdir("..")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(wd)

	out := e2e.Run(t, "run", "file_write_all")

	require.Equal(
		t,
		"",
		strings.TrimSuffix(out, "\n"),
	)

	// Check file contents.
	const filename = "file_writer_example.txt"

	want, err := os.ReadFile(filename)
	require.NoError(t, err, out)
	require.Equal(
		t,
		"Hello, io.WriteAll!",
		string(want),
	)

	// Remove file output.
	os.Remove(filename)
}
