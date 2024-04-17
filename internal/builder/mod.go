package builder

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
)

func (p Builder) LoadModuleByPath(
	ctx context.Context,
	wd string,
) (compiler.RawModule, error) {
	manifest, modRootPath, err := p.getNearestManifest(wd)
	if err != nil {
		return compiler.RawModule{}, fmt.Errorf("retrieve manifest: %w", err)
	}

	pkgs := map[string]compiler.RawPackage{}
	if err := retrieveSourceCode(modRootPath, pkgs); err != nil {
		return compiler.RawModule{}, fmt.Errorf("walk: %w", err)
	}

	return compiler.RawModule{
		Manifest: manifest,
		Packages: pkgs,
	}, nil
}

// retrieveSourceCode recursively walks the given tree and fills given pkgs with neva files
func retrieveSourceCode(rootPath string, pkgs map[string]compiler.RawPackage) error {
	fsys := os.DirFS(rootPath)
	return fs.WalkDir(fsys, ".", func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("filepath walk: %s: %w", filePath, err)
		}

		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(d.Name())
		if ext != ".neva" {
			return nil
		}

		file, err := fsys.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		bb, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		pkgName := getPkgName(rootPath, filePath)
		if _, ok := pkgs[pkgName]; !ok {
			pkgs[pkgName] = compiler.RawPackage{}
		}

		fileName := strings.TrimSuffix(d.Name(), ext)
		pkgs[pkgName][fileName] = bb

		return nil
	})
}

func getPkgName(rootPath, filePath string) string {
	dirPath := filepath.Dir(filePath)
	if dirPath == rootPath { // current directory is root directory
		ss := strings.Split(dirPath, "/")
		return ss[len(ss)-1] // return last element
	}
	return strings.TrimPrefix(dirPath, rootPath+"/")
}
