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
	out, _ := e2e.Run(t, []string{"run", "."})

	require.Equal(
		t,
		"",
		strings.TrimSuffix(out, "\n"),
	)

	// Check file contents.
	repoRoot := e2e.FindRepoRoot(t)
	filename := filepath.Join(repoRoot, "examples", "file_write_all", "file_writer_example.txt")

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
