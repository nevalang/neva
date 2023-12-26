package pkgmanager

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/pkg/sourcecode"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

func (p Manager) buildModule(ctx context.Context, workdir string) (compiler.RawModule, error) {
	manifest, err := p.retrieveManifest(workdir)
	if err != nil {
		return compiler.RawModule{}, nil
	}

	pkgs := map[string]compiler.RawPackage{}
	if err := walkTree(workdir, pkgs); err != nil {
		return compiler.RawModule{}, fmt.Errorf("walk: %w", err)
	}

	return compiler.RawModule{
		Manifest: manifest,
		Packages: pkgs,
	}, nil
}

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

func walkTree(rootPath string, pkgs map[string]compiler.RawPackage) error {
	if err := filepath.Walk(rootPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("filepath walk: %s: %w", filePath, err)
		}

		ext := filepath.Ext(info.Name())
		if ext != ".neva" {
			return nil
		}

		bb, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		pkgName := getPkgName(rootPath, filePath)
		if _, ok := pkgs[pkgName]; !ok {
			pkgs[pkgName] = compiler.RawPackage{}
		}

		fileName := strings.TrimSuffix(info.Name(), ext)
		pkgs[pkgName][fileName] = bb

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func getPkgName(rootPath, filePath string) string {
	dirPath := filepath.Dir(filePath)
	if dirPath == rootPath { // current directory is root directory
		ss := strings.Split(dirPath, "/")
		return ss[len(ss)-1] // return last element
	}
	return strings.TrimPrefix(dirPath, rootPath+"/")
}
