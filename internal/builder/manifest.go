package builder

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	yaml "gopkg.in/yaml.v3"

	"github.com/nevalang/neva/internal/compiler/sourcecode"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

func (p Builder) getNearestManifest(filesys fs.FS, wd string) (src.ModuleManifest, error) {
	rawNearest, err := lookupManifestFile(filesys, wd, 0)
	if err != nil {
		return sourcecode.ModuleManifest{}, fmt.Errorf("read manifest yaml: %w", err)
	}

	parsedNearest, err := p.manifestParser.ParseManifest(rawNearest)
	if err != nil {
		return sourcecode.ModuleManifest{}, fmt.Errorf("parse manifest: %w", err)
	}

	return parsedNearest, nil
}

func lookupManifestFile(filesys fs.FS, wd string, iteration int) ([]byte, error) {
	if iteration > 10 {
		return nil, errors.New("manifest file not found")
	}

	found, err := readManifestFromDir(filesys, wd)
	if err == nil {
		return found, nil
	}

	if !errors.Is(err, fs.ErrInvalid) {
		return nil, err
	}

	return lookupManifestFile(
		filesys,
		path.Dir(wd),
		iteration+1,
	)
}

func readManifestFromDir(filesys fs.FS, wd string) ([]byte, error) {
	raw, err := fs.ReadFile(filesys, path.Join(wd, "neva.yaml"))
	if err == nil {
		return raw, nil
	}

	return fs.ReadFile(filesys, path.Join(wd, "neva.yml"))
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
