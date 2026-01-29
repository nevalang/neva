package native

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
	tmpGoModuleDir, err := os.MkdirTemp(dst, "neva_build_")
	if err != nil {
		return fmt.Errorf("create temporary build directory: %w", err)
	}

	if err := b.golang.EmitExecutable(tmpGoModuleDir, prog, trace); err != nil {
		return fmt.Errorf("emit executable: %w", err)
	}

	if err := b.buildExecutable(tmpGoModuleDir, dst); err != nil {
		return fmt.Errorf("build executable: %w", err)
	}

	if err := os.RemoveAll(tmpGoModuleDir); err != nil {
		return fmt.Errorf("remove gomodule: %w", err)
	}

	return nil
}

func (b Backend) EmitLibrary(dst string, exports []compiler.LibraryExport, trace bool) error {
	return fmt.Errorf("library mode not implemented for native backend")
}

func (b Backend) buildExecutable(gomodule, output string) error {
	fileName := "output"
	if os.Getenv("GOOS") == "windows" { // either we're on windows or we're cross-compiling
		fileName += ".exe"
	}

	// #nosec G204 -- command args are constructed internally from known values
	cmd := exec.Command(
		"go",
		"build",
		"-ldflags", "-s -w", // strip debug information
		"-o",
		filepath.Join(output, fileName),
		".",
	)
	cmd.Dir = gomodule
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func NewBackend(golangBackend golang.Backend) Backend {
	return Backend{
		golang: golangBackend,
	}
}
