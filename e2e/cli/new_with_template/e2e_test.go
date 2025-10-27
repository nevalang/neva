package test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/require"
)

func TestNewFromTemplate(t *testing.T) {
	templateDir := t.TempDir()

	repo, err := git.PlainInit(templateDir, false)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "neva.yml"), []byte("neva: 0.0.1"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(templateDir, "src"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "src", "main.neva"), []byte("def Main() {}"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(templateDir, "nested"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "nested", "file.txt"), []byte("data"), 0o755))

	worktree, err := repo.Worktree()
	require.NoError(t, err)
	_, err = worktree.Add("neva.yml")
	require.NoError(t, err)
	_, err = worktree.Add("src/main.neva")
	require.NoError(t, err)
	_, err = worktree.Add("nested/file.txt")
	require.NoError(t, err)
	_, err = worktree.Commit("initial", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)

	workdir := t.TempDir()
	destination := filepath.Join(workdir, "module")

	cmd := exec.Command("neva", "new", destination, "--template", templateDir)
	cmd.Dir = workdir

	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))

	_, err = os.Stat(filepath.Join(destination, "neva.yml"))
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(destination, "src", "main.neva"))
	require.NoError(t, err)
	info, err := os.Stat(filepath.Join(destination, "nested", "file.txt"))
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o755), info.Mode().Perm())

	_, err = os.Stat(filepath.Join(destination, ".git"))
	require.ErrorIs(t, err, os.ErrNotExist)
}
