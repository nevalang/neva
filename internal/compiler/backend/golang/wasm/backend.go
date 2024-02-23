package wasm

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nevalang/neva/internal/compiler/backend/golang"
	"github.com/nevalang/neva/internal/runtime/ir"
)

type Backend struct {
	golang golang.Backend
}

func (b Backend) Emit(dst string, prog *ir.Program) error {
	tmpGoProj := dst + "/tmp"
	if err := b.golang.Emit(tmpGoProj, prog); err != nil {
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

// TODO handle the whole pipeline including html and js glue generation.

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
