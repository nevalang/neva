package parser

import (
	"fmt"

	yaml "gopkg.in/yaml.v3"

	src "github.com/nevalang/neva/pkg/sourcecode"
)

func (p Parser) ParseManifest(raw []byte) (src.ModuleManifest, error) {
	var manifest Manifest
	if err := yaml.Unmarshal(raw, &manifest); err != nil {
		return src.ModuleManifest{}, fmt.Errorf("yaml unmarshal: %w", err)
	}
	return manifestToSourceCode(manifest), nil
}

type Manifest struct {
	LanguageVersion string      `yaml:"neva"`
	Deps            []ModuleRef `yaml:"deps"`
}

type ModuleRef struct {
	Path    string `yaml:"path"`
	Version string `yaml:"version"`
	Alias   string `yaml:"alias"`
}

func manifestToSourceCode(manifest Manifest) src.ModuleManifest {
	deps := make(map[string]src.ModuleRef, len(manifest.Deps))
	for _, dep := range manifest.Deps {
		var k string
		if dep.Alias != "" {
			k = dep.Alias
		} else {
			k = dep.Path
		}
		deps[k] = src.ModuleRef{
			Path:    dep.Path,
			Version: dep.Version,
		}
	}
	return src.ModuleManifest{
		LanguageVersion: manifest.LanguageVersion,
		Deps:            deps,
	}
}
