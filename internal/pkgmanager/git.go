package pkgmanager

import (
	"fmt"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/nevalang/neva/pkg/sourcecode"
)

func (p Manager) downloadDep(depModRef sourcecode.ModuleRef) (string, error) {
	fsPath := fmt.Sprintf("%s/%s_%s", p.thirdPartyLocation, depModRef.Path, depModRef.Version)

	_, err := os.Stat(fsPath)
	if err == nil {
		return fsPath, nil
	}

	if !os.IsNotExist(err) {
		return "", fmt.Errorf("os stat: %w", err)
	}

	if _, err := git.PlainClone(fsPath, false, &git.CloneOptions{
		URL:           "https://" + depModRef.Path,
		ReferenceName: plumbing.NewTagReferenceName(depModRef.Version),
	}); err != nil {
		return "", err
	}

	return fsPath, nil
}
