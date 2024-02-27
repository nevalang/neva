package builder

import (
	"fmt"
	"io/fs"

	"github.com/nevalang/neva/pkg/sourcecode"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

func (p Builder) retrieveManifest(workdir fs.FS) (src.ModuleManifest, error) {
	rawManifest, err := readManifestYaml(workdir)
	if err != nil {
		return sourcecode.ModuleManifest{}, fmt.Errorf("read manifest yaml: %w", err)
	}

	manifest, err := p.manifestParser.ParseManifest(rawManifest)
	if err != nil {
		return sourcecode.ModuleManifest{}, fmt.Errorf("parse manifest: %w", err)
	}

	return manifest, nil
}

func readManifestYaml(workdir fs.FS) ([]byte, error) {
	rawManifest, err := fs.ReadFile(workdir, "neva.yml")
	if err == nil {
		return rawManifest, nil
	}

	rawManifest, err = fs.ReadFile(workdir, "neva.yaml")
	if err != nil {
		return nil, fmt.Errorf("fs read file: %w", err)
	}

	return rawManifest, nil
}
