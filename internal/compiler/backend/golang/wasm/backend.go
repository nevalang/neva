package wasm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nevalang/neva/internal/compiler"
	backendgolang "github.com/nevalang/neva/internal/compiler/backend/golang"
	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/nevalang/neva/pkg/golang"
)

type Backend struct {
	golang backendgolang.Backend
}

func (b Backend) EmitExecutable(dst string, prog *ir.Program, trace bool) error {
	tmpGoProj := dst + "/tmp"
	if err := b.golang.EmitExecutable(tmpGoProj, prog, trace); err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return err
	}
	if err := buildWASM(tmpGoProj, dst); err != nil {
		return err
	}
	//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	return os.RemoveAll(tmpGoProj)
}

func (b Backend) EmitLibrary(dst string, exports []compiler.LibraryExport, trace bool) error {
	//nolint:perfsprint // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	return fmt.Errorf("library mode not implemented for wasm backend")
}

func buildWASM(src, dst string) error {
	outputPath := filepath.Join(dst, "output")
	if err := os.Chdir(src); err != nil {
		//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
		return err
	}
	// #nosec G204 -- command args are constructed internally from known values
	//nolint:noctx // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	cmd := exec.Command(
		"go",
		golang.ReleaseBuildArgs(outputPath+".wasm", src)...,
	)
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	//nolint:wrapcheck // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
	return cmd.Run()
}

func NewBackend(golangBackend backendgolang.Backend) Backend {
	return Backend{
		golang: golangBackend,
	}
}
