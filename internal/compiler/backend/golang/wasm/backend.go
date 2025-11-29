package wasm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/backend/golang"
	"github.com/nevalang/neva/internal/compiler/ir"
)

type Backend struct {
	golang golang.Backend
}

func (b Backend) EmitExecutable(dst string, prog *ir.Program, trace bool) error {
	tmpGoProj := dst + "/tmp"
	if err := b.golang.EmitExecutable(tmpGoProj, prog, trace); err != nil {
		return err
	}
	if err := buildWASM(tmpGoProj, dst); err != nil {
		return err
	}
	if err := os.RemoveAll(tmpGoProj); err != nil {
		return err
	}
	return nil
}

func (b Backend) EmitLibrary(dst string, exports []compiler.LibraryExport, trace bool) error {
	return fmt.Errorf("library mode not implemented for wasm backend")
}

func buildWASM(src, dst string) error {
	outputPath := filepath.Join(dst, "output")
	if err := os.Chdir(src); err != nil {
		return err
	}
	cmd := exec.Command(
		"go",
		"build",
		"-ldflags", "-s -w", // for optimization
		"-o", outputPath+".wasm",
		src,
	)
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func NewBackend(golangBackend golang.Backend) Backend {
	return Backend{
		golang: golangBackend,
	}
}
