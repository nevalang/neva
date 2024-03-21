package builder

import (
	"errors"
	"fmt"
	"os"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

func (b Builder) Get(workdir, path, version string) (string, error) {
	ref := src.ModuleRef{
		Path:    path,
		Version: version,
	}

	downloadPath, actualVersion, err := b.downloadDep(ref)
	if err != nil {
		return "", err
	}

	manifest, err := b.retrieveManifest(os.DirFS(workdir))
	if err != nil {
		return "", fmt.Errorf("Retrieve manifest: %w", err)
	}

	existing, ok := manifest.Deps[path]
	if ok && existing.Version != actualVersion {
		return "", errors.New(
			"Several versions of the same dependency not yet supported.",
		)
	}

	manifest.Deps[path] = src.ModuleRef{
		Path:    path,
		Version: actualVersion,
	}

	if err := b.writeManifest(manifest, workdir); err != nil {
		return "", err
	}

	return downloadPath, nil
}
