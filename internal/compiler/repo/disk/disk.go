package disk

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/pkg/ir"
	"google.golang.org/protobuf/proto"
)

type Repository struct {
	stdPath string
}

func (r Repository) ByPath(ctx context.Context, pathToMainPkg string) (map[string]compiler.RawPackage, error) {
	mainPkgFiles := map[string][]byte{}
	if err := readAllFilesInDir(pathToMainPkg, mainPkgFiles); err != nil {
		return nil, err
	}

	stdFiles := map[string][]byte{}
	if err := readAllFilesInDir(r.stdPath, stdFiles); err != nil {
		return nil, err
	}

	return map[string]compiler.RawPackage{
		"main": mainPkgFiles,
		"std":  stdFiles,
	}, nil
}

func (r Repository) Save(ctx context.Context, path string, prog *ir.Program) error {
	bb, err := proto.Marshal(prog)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bb, 0644) //nolint:gofumpt
}

func readAllFilesInDir(path string, files map[string][]byte) error {
	if err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// err := readAllFilesInDir(filePath, files)
			// if err != nil {
			// 	return err
			// }
			return nil
		}
		ext := filepath.Ext(info.Name())
		if ext != ".neva" {
			return nil
		}
		bb, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		files[strings.TrimSuffix(info.Name(), ext)] = bb
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func MustNew(stdPkgPath string) Repository {
	return Repository{
		stdPath: stdPkgPath,
	}
}
