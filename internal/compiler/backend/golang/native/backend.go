package native

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nevalang/neva/internal/compiler/backend/golang"
	"github.com/nevalang/neva/internal/compiler/ir"
)

type Backend struct {
	golang golang.Backend
}

func (b Backend) Emit(output string, prog *ir.Program, trace bool) error {
	tmpGoModuleDir := output + "/tmp"
	if err := b.golang.Emit(tmpGoModuleDir, prog, trace); err != nil {
		return fmt.Errorf("emit: %w", err)
	}
	if err := b.buildExecutable(tmpGoModuleDir, output); err != nil {
		return fmt.Errorf("build executable: %w", err)
	}
	if err := os.RemoveAll(tmpGoModuleDir); err != nil {
		return fmt.Errorf("remove gomodule: %w", err)
	}
	return nil
}

func (b Backend) buildExecutable(gomodule, output string) error {
	// remember current working directory to change back to it later
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	// we need to be inside go module to run `go build` command
	if err := os.Chdir(gomodule); err != nil {
		return fmt.Errorf("change directory to gomodule: %w", err)
	}

	// change back to original wd or neva programs
	// that interact with fs via relative paths will fail
	defer func() {
		if err := os.Chdir(wd); err != nil {
			panic(err)
		}
	}()

	fileName := "output"
	if os.Getenv("GOOS") == "windows" { // either we're on windows or we're cross-compiling
		fileName += ".exe"
	}

	cmd := exec.Command(
		"go",
		"build",
		"-ldflags", "-s -w", // strip debug information
		"-o",
		filepath.Join(output, fileName),
		gomodule,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func NewBackend(golangBackend golang.Backend) Backend {
	return Backend{
		golang: golangBackend,
	}
}
