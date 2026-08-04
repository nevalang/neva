package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

const unformattedSource = "def Main(start any) (stop any) {\n:start->:stop\n}\n"

const formattedSource = "def Main(start any) (stop any) {\n\t:start -> :stop\n}\n"

// TestFmt verifies the installed command's standard-input, write, check, and
// command-line error contracts through a real subprocess.
func TestFmt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "main.neva")
	require.NoError(t, os.WriteFile(path, []byte(unformattedSource), 0o600))

	stdout, _ := e2e.Run(t, []string{"fmt"}, e2e.WithDir(dir), e2e.WithStdin(unformattedSource))
	require.Equal(t, formattedSource, stdout)

	stdout, _ = e2e.Run(t, []string{"fmt", "-check", path}, e2e.WithDir(dir), e2e.WithCode(1))
	require.Equal(t, path+"\n", stdout)

	e2e.Run(t, []string{"fmt", "-w", path}, e2e.WithDir(dir))
	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, formattedSource, string(contents))

	stdout, _ = e2e.Run(t, []string{"fmt", "-check", path}, e2e.WithDir(dir))
	require.Empty(t, stdout)

	e2e.Run(t, []string{"fmt", "-w", "-d", path}, e2e.WithDir(dir), e2e.WithCode(2))
}
