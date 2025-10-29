package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFromRemoteTemplate(t *testing.T) {
	workdir := t.TempDir()
	destination := filepath.Join(workdir, "module")

	cmd := exec.Command("neva", "new", destination, "--template", "github.com/nevalang/x")
	cmd.Dir = workdir

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	// verify that files from the template were copied
	_, err = os.Stat(filepath.Join(destination, "neva.yml"))
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(destination, "src"))
	require.NoError(t, err)

	// verify that .git directory was not copied
	_, err = os.Stat(filepath.Join(destination, ".git"))
	require.ErrorIs(t, err, os.ErrNotExist)
}
