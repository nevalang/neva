package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nevalang/neva/pkg/e2e"
	"github.com/stretchr/testify/require"
)

func TestNewFromRemoteTemplate(t *testing.T) {
	workdir := t.TempDir()
	destination := filepath.Join(workdir, "module")

	e2e.RunInDir(t, workdir, "new", destination, "--template", "github.com/nevalang/x")

	// verify that files from the template were copied
	_, err := os.Stat(filepath.Join(destination, "neva.yml"))
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(destination, "src"))
	require.NoError(t, err)

	// verify that .git directory was not copied
	_, err = os.Stat(filepath.Join(destination, ".git"))
	require.ErrorIs(t, err, os.ErrNotExist)
}
