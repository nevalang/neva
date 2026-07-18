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
	repoRoot := e2e.FindRepoRoot(t)
	exampleDir := filepath.Join(repoRoot, "examples", "file_handles")

	workdir, err := os.MkdirTemp(exampleDir, ".file-handles-e2e-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(workdir))
	})

	src := filepath.Join(exampleDir, "main.neva")
	dst := filepath.Join(workdir, "main.neva")
	data, err := os.ReadFile(src)
	require.NoError(t, err)
	// #nosec G703 -- dst is inside a test-created temporary directory.
	require.NoError(t, os.WriteFile(dst, data, 0o600))

	out, _ := e2e.Run(t, []string{"run", "."}, e2e.WithDir(workdir))

	require.Equal(t, "Hello, io.File!", strings.TrimSuffix(out, "\n"))

	filename := filepath.Join(workdir, "file_handles_example.txt")
	contents, err := os.ReadFile(filename)
	require.NoError(t, err, out)
	require.Equal(t, "Hello, io.File!", string(contents))
}
