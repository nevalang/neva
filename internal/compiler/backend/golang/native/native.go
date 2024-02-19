package native

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nevalang/neva/internal/compiler/backend/golang"
	"github.com/nevalang/neva/pkg/ir"
)

type Backend struct {
	golang golang.Backend
}

func (b Backend) Emit(dst string, prog *ir.Program) error {
	gomod := dst + "/tmp"
	if err := b.golang.Emit(gomod, prog); err != nil {
		return err
	}
	if err := buildExecutable(gomod, dst); err != nil {
		return err
	}
	if err := os.RemoveAll(gomod); err != nil {
		return err
	}
	return nil
}

func buildExecutable(src, dst string) error {
	outputPath := filepath.Join(dst, "output")
	if err := os.Chdir(src); err != nil {
		return err
	}
	cmd := exec.Command("go", "build", "-o", outputPath, src)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func NewBackend(golangBackend golang.Backend) Backend {
	return Backend{
		golang: golangBackend,
	}
}
