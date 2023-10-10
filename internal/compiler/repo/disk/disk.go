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
	stdLibPath string
}

func (r Repository) ByPath(ctx context.Context, pathToMainPkg string) (map[string]compiler.RawPackage, error) {
	mainPkgFiles := map[string][]byte{}
	if err := readAllFilesInDir(pathToMainPkg, mainPkgFiles); err != nil {
		return nil, err
	}

	prog := map[string]compiler.RawPackage{
		"main": mainPkgFiles,
	}

	stdTmpFiles := map[string][]byte{}
	if err := readAllFilesInDir(r.stdLibPath+"/tmp", stdTmpFiles); err != nil {
		return nil, err
	}

	stdBuiltinFiles := map[string][]byte{}
	if err := readAllFilesInDir(r.stdLibPath+"/builtin", stdBuiltinFiles); err != nil {
		return nil, err
	}

	prog["std/tmp"] = stdTmpFiles
	prog["std/builtin"] = stdBuiltinFiles

	return prog, nil
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
		stdLibPath: stdPkgPath,
	}
}
