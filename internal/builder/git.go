package builder

import (
	"fmt"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/nevalang/neva/internal/compiler/sourcecode"
)

func (b Builder) downloadDep(dep sourcecode.ModuleRef) (string, error) {
	path := fmt.Sprintf("%s/%s_%s", b.thirdPartyLocation, dep.Name, dep.Version)

	_, err := os.Stat(path)
	if err == nil {
		return path, nil
	}

	if !os.IsNotExist(err) {
		return "", fmt.Errorf("os stat: %w", err)
	}

	if _, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:           fmt.Sprintf("https://%s", dep.Name),
		ReferenceName: plumbing.NewTagReferenceName(dep.Version),
	}); err != nil {
		return "", err
	}

	return path, nil
}
