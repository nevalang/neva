package builder

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/ast/core"
)

func (b Builder) Get(wd, path, version string) (string, error) {
	ref := core.ModuleRef{
		Path:    path,
		Version: version,
	}

	release, err := acquireLockFile()
	if err != nil {
		return "", fmt.Errorf("failed to acquire lock file: %w", err)
	}
	defer release()

	downloadPath, actualVersion, err := b.downloadDep(ref)
	if err != nil {
		return "", err
	}

	manifest, _, err := b.getNearestManifest(wd)
	if err != nil {
		return "", fmt.Errorf("Retrieve manifest: %w", err)
	}

	existing, ok := manifest.Deps[path]
	if ok && existing.Version != actualVersion {
		return "", errors.New(
			"Several versions of the same dependency not yet supported.",
		)
	}

	manifest.Deps[path] = core.ModuleRef{
		Path:    path,
		Version: actualVersion,
	}

	if err := b.writeManifest(manifest, wd); err != nil {
		return "", err
	}

	return downloadPath, nil
}
