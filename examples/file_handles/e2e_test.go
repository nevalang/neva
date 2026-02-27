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

	require.Equal(t, "Hello, io.File!", strings.TrimSuffix(out, "\n"))

	repoRoot := e2e.FindRepoRoot(t)
	filename := filepath.Join(repoRoot, "examples", "file_handles", "file_handles_example.txt")
	contents, err := os.ReadFile(filename)
	require.NoError(t, err, out)
	require.Equal(t, "Hello, io.File!", string(contents))

	require.NoError(t, os.Remove(filename))
}
