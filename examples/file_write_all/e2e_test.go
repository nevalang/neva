package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	out := e2e.RunExample(t, "file_write_all")

	require.Equal(
		t,
		"",
		strings.TrimSuffix(out, "\n"),
	)

	// Check file contents.
	filename := filepath.Join(e2e.ExamplesDir(t), "file_writer_example.txt")

	want, err := os.ReadFile(filename)
	require.NoError(t, err, out)
	require.Equal(
		t,
		"Hello, io.WriteAll!",
		string(want),
	)

	// Remove file output.
	_ = os.Remove(filename)
}
