package disk

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
)

type Builder struct {
	stdLibPath string
}

func (r Builder) Build(ctx context.Context, workdir string) (map[string]compiler.RawPackage, error) {
	prog := map[string]compiler.RawPackage{}

	// read all packages in working directory recursively
	if err := walk(workdir, prog); err != nil {
		return nil, fmt.Errorf("walk: %w", err)
	}

	// read all packages in stdlib directory recursively
	if err := walk(r.stdLibPath, prog); err != nil {
		return nil, fmt.Errorf("walk: %w", err)
	}

	return prog, nil
}

func walk(path string, prog map[string]compiler.RawPackage) error {
	if err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("filepath walk: %s: %w", filePath, err)
		}

		if info.IsDir() {
			return walk(filePath, prog)
		}

		ext := filepath.Ext(info.Name())
		if ext != ".neva" {
			return nil
		}

		bb, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		if _, ok := prog[path]; !ok {
			prog[path] = compiler.RawPackage{}
		}

		fileName := strings.TrimSuffix(info.Name(), ext)
		prog[path][fileName] = bb

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func MustNew(stdPkgPath string) Builder {
	return Builder{
		stdLibPath: stdPkgPath,
	}
}
