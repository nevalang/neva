package cli

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	nevaGit "github.com/nevalang/neva/pkg/git"
	"github.com/stretchr/testify/require"
)

func TestScaffoldFromTemplate(t *testing.T) {
	t.Parallel()

	// Bootstraps a local git repository to assert the scaffold copies files
	// verbatim (preserving permissions) while dropping VCS metadata.
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

	destination := filepath.Join(t.TempDir(), "module")
	spec, err := nevaGit.ParseRepoSpec(templateDir)
	require.NoError(t, err)
	require.NoError(t, scaffoldFromTemplate(destination, spec))

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
