package pkgmanager

import (
	"fmt"
	"os"

	"github.com/nevalang/neva/pkg/sourcecode"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

func (p Manager) retrieveManifest(workdir string) (src.ModuleManifest, error) {
	rawManifest, err := readManifestYaml(workdir)
	if err != nil {
		return sourcecode.ModuleManifest{}, fmt.Errorf("read manifest yaml: %w", err)
	}

	manifest, err := p.parser.ParseManifest(rawManifest)
	if err != nil {
		return sourcecode.ModuleManifest{}, fmt.Errorf("parse manifest: %w", err)
	}

	return manifest, nil
}

func readManifestYaml(workdir string) ([]byte, error) {
	rawManifest, err := os.ReadFile(workdir + "/neva.yml")
	if err == nil {
		return rawManifest, nil
	}

	rawManifest, err = os.ReadFile(workdir + "/neva.yaml")
	if err != nil {
		return nil, fmt.Errorf("os read file: %w", err)
	}

	return rawManifest, nil
}
