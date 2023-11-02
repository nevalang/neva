package builder

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/src"
)

type Builder struct {
	stdLibPath string
}

func (r Builder) BuildModule(ctx context.Context, workdir string) (compiler.RawModule, error) {
	prog := map[string]compiler.RawPackage{}

	// read all packages in working directory recursively
	if err := walk(workdir, prog); err != nil {
		return compiler.RawModule{}, fmt.Errorf("walk: %w", err)
	}

	// read all packages in stdlib directory recursively
	if err := walk(r.stdLibPath, prog); err != nil {
		return compiler.RawModule{}, fmt.Errorf("walk: %w", err)
	}

	return compiler.RawModule{
		Manifest: src.Manifest{}, // TODO
		Packages: map[string]compiler.RawPackage{},
	}, nil
}

// walk recursively traverses the directory assuming that is a neva module (set of packages with source code files).
func walk(rootPath string, prog map[string]compiler.RawPackage) error {
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
		if _, ok := prog[pkgName]; !ok {
			prog[pkgName] = compiler.RawPackage{}
		}

		fileName := strings.TrimSuffix(info.Name(), ext)
		prog[pkgName][fileName] = bb

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

func MustNew(stdPkgPath string) Builder {
	return Builder{
		stdLibPath: stdPkgPath,
	}
}
