package builder

import (
	"fmt"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/nevalang/neva/internal/compiler/sourcecode"
)

// downloadDep returns path where it downloaded dependency
// and its downloaded version in case version wasn't specified.
func (p Builder) downloadDep(depModRef sourcecode.ModuleRef) (string, string, error) {
	fsPath := fmt.Sprintf(
		"%s/%s_%s",
		p.thirdPartyPath,
		depModRef.Path,
		depModRef.Version,
	)

	_, err := os.Stat(fsPath)
	if err == nil {
		return fsPath, depModRef.Version, nil
	}

	if !os.IsNotExist(err) {
		return "", "", fmt.Errorf("os stat: %w", err)
	}

	var ref plumbing.ReferenceName
	if depModRef.Version != "" {
		ref = plumbing.NewTagReferenceName(depModRef.Version)
	}

	repo, err := git.PlainClone(fsPath, false, &git.CloneOptions{
		URL:           "https://" + depModRef.Path,
		ReferenceName: ref,
	})
	if err != nil {
		return "", "", err
	}

	if ref != "" {
		return fsPath, "", nil
	}

	latestTagHash, tagName, err := getLatestTagHash(repo)
	if err != nil {
		return "", "", err
	}

	tree, err := repo.Worktree()
	if err != nil {
		return "", "", err
	}

	if err := tree.Checkout(&git.CheckoutOptions{
		Hash: latestTagHash,
	}); err != nil {
		return "", "", err
	}

	// Append the latest tag to the directory name
	newFsPath := fmt.Sprintf("%s%s", fsPath, tagName)

	// Remove the directory if it already exists
	if _, err := os.Stat(newFsPath); err == nil {
		if err := os.RemoveAll(newFsPath); err != nil {
			return "", "", fmt.Errorf("os.RemoveAll: %w", err)
		}
	}

	// Finally rename new directory
	if err := os.Rename(fsPath, newFsPath); err != nil {
		return "", "", fmt.Errorf("os rename: %w", err)
	}

	return newFsPath, tagName, nil
}

func getLatestTagHash(repository *git.Repository) (plumbing.Hash, string, error) {
	tagRefs, err := repository.Tags()
	if err != nil {
		return plumbing.Hash{}, "", err
	}

	var (
		hash plumbing.Hash
		name string
	)
	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		hash = tagRef.Hash()
		name = tagRef.Name().String()
		return nil
	})
	if err != nil {
		return plumbing.Hash{}, "", err
	}

	return hash, strings.Split(name, "/")[2], nil
}
