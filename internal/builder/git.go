package builder

import (
	"fmt"
	"os"
	"strings"

	gitlib "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/nevalang/neva/pkg/core"
	nevaGit "github.com/nevalang/neva/pkg/git"
)

// downloadDep returns path where it downloaded dependency
// and its downloaded version in case version wasn't specified.
func (p Builder) downloadDep(depModRef core.ModuleRef) (string, string, error) {
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

	spec, err := nevaGit.ParseRepoSpec(depModRef.Path)
	if err != nil {
		return "", "", fmt.Errorf("parse module repository: %w", err)
	}

	repo, err := gitlib.PlainClone(fsPath, false, &gitlib.CloneOptions{
		URL:           spec.CloneURL(),
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

	if err := nevaGit.Checkout(repo, latestTagHash.String()); err != nil {
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

func getLatestTagHash(repository *gitlib.Repository) (plumbing.Hash, string, error) {
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
