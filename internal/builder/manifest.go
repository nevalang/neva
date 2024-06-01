package builder

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"

	"github.com/nevalang/neva/internal/compiler/sourcecode"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

func (p Builder) getNearestManifest(wd string) (src.ModuleManifest, string, error) {
	rawNearest, path, err := lookupManifestFile(wd, 0)
	if err != nil {
		return sourcecode.ModuleManifest{}, "", fmt.Errorf("read manifest yaml: %w", err)
	}

	parsedNearest, err := p.manifestParser.ParseManifest(rawNearest)
	if err != nil {
		return sourcecode.ModuleManifest{}, "", fmt.Errorf("parse manifest: %w", err)
	}

	return parsedNearest, path, nil
}

func lookupManifestFile(wd string, iteration int) ([]byte, string, error) {
	if iteration > 10 {
		return nil, "", errors.New("manifest file not found in 10 nearest levels up to where cli executed")
	}

	found, err := readManifestFromDir(wd)
	if err == nil {
		return found, wd, nil
	}

	if !errors.Is(err, fs.ErrInvalid) &&
		!errors.Is(err, os.ErrNotExist) {
		return nil, "", err
	}

	return lookupManifestFile(
		filepath.Dir(wd),
		iteration+1,
	)
}

func readManifestFromDir(wd string) ([]byte, error) {
	raw, err := os.ReadFile(filepath.Join(wd, "neva.yaml"))
	if err == nil {
		return raw, nil
	}
	return os.ReadFile(filepath.Join(wd, "neva.yml"))
}

func (b Builder) writeManifest(manifest src.ModuleManifest, workdir string) error {
	manifestData, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}

	manifestFileName := "neva.yaml"
	if _, err := os.Stat(filepath.Join(workdir, "neva.yml")); err == nil {
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
