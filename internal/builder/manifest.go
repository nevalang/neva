package builder

import (
	"fmt"
	"io/fs"
	"os"

	yaml "gopkg.in/yaml.v3"

	"github.com/nevalang/neva/internal/compiler/sourcecode"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
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
		files, readDirErr := fs.ReadDir(workdir, ".")
		if readDirErr != nil {
			return nil, fmt.Errorf("fs read dir: %w", readDirErr)
		}
		for _, file := range files {
			fmt.Println(file.Name())
		}

		return nil, fmt.Errorf("fs read file: %w", err)
	}

	return rawManifest, nil
}

func (b Builder) writeManifest(manifest src.ModuleManifest, workdir string) error {
	manifestData, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}

	manifestFileName := "neva.yaml"
	if _, err := os.Stat(workdir + "/neva.yml"); err == nil {
		manifestFileName = "neva.yml"
	}

	manifestPath := workdir + "/" + manifestFileName
	file, err := os.OpenFile(manifestPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(manifestData)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
