package disk

import (
	"context"
	"os"
	"path/filepath"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/pkg/ir"
	"google.golang.org/protobuf/proto"
)

type Repository struct{}

func (r Repository) ByPath(ctx context.Context, pathToMainPkg string) (map[string]compiler.RawPackage, error) {
	mainPkgFiles, err := readAllFilesInDir(pathToMainPkg)
	if err != nil {
		return nil, err
	}
	return map[string]compiler.RawPackage{
		"main": mainPkgFiles,
	}, nil
}

func (r Repository) Save(ctx context.Context, path string, prog *ir.Program) error {
	bb, err := proto.Marshal(prog)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bb, 0644)
}

func readAllFilesInDir(path string) (map[string][]byte, error) {
	files := make(map[string][]byte)

	if err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileBytes, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}
			files[filePath] = fileBytes
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}
func MustNew() Repository {
	return Repository{}
}
