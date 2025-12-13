package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	cli "github.com/urfave/cli/v2"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/backend/golang"
	"github.com/nevalang/neva/internal/compiler/backend/golang/native"
	"github.com/nevalang/neva/internal/compiler/desugarer"
)

func newInstallCmd(
	workdir string,
	bldr builder.Builder,
	parser compiler.Parser,
	desugarer desugarer.Desugarer,
	analyzer compiler.Analyzer,
	irgen compiler.Irgen,
) *cli.Command {
	return &cli.Command{
		Name:      "install",
		Usage:     "Build and install neva program",
		Args:      true,
		ArgsUsage: "Provide path to main package",
		Action: func(cliCtx *cli.Context) error {
			mainPkg, err := mainPkgPathFromArgs(cliCtx)
			if err != nil {
				return err
			}

			// Resolve absolute path to package
			absPkg := mainPkg
			if !filepath.IsAbs(absPkg) {
				absPkg = filepath.Join(workdir, absPkg)
			}

			absPkg = filepath.Clean(absPkg)
			pkgBase := absPkg

			pkgInfo, err := os.Stat(absPkg)
			if err != nil {
				return fmt.Errorf("stat package: %w", err)
			}

			// If path points to a file, use its directory as base
			if !pkgInfo.IsDir() {
				pkgBase = filepath.Dir(absPkg)
			}

			// Resolve actual package path (handles module root vs src subdirectory)
			compilePkg, err := resolveMainPkgPath(absPkg)
			if err != nil {
				return err
			}

			// Binary name is derived from package directory name
			binName := filepath.Base(pkgBase)
			if binName == string(filepath.Separator) || binName == "." || binName == "" {
				return fmt.Errorf("cannot determine binary name for %s", mainPkg)
			}

			// Ensure GOOS/GOARCH match runtime to build for current platform
			prevGOOS := os.Getenv("GOOS")
			prevGOARCH := os.Getenv("GOARCH")

			if err := os.Setenv("GOOS", runtime.GOOS); err != nil {
				return fmt.Errorf("set GOOS: %w", err)
			}
			if err := os.Setenv("GOARCH", runtime.GOARCH); err != nil {
				return fmt.Errorf("set GOARCH: %w", err)
			}
			defer func() {
				if err := os.Setenv("GOOS", prevGOOS); err != nil {
					panic(err)
				}
				if err := os.Setenv("GOARCH", prevGOARCH); err != nil {
					panic(err)
				}
			}()

			// Compile to temporary directory first
			tempDir, err := os.MkdirTemp("", "neva-install-*")
			if err != nil {
				return fmt.Errorf("create temp dir: %w", err)
			}
			defer os.RemoveAll(tempDir)

			compilerToNative := compiler.New(
				bldr,
				parser,
				&desugarer,
				analyzer,
				irgen,
				native.NewBackend(golang.NewBackend("")),
			)

			if _, err := compilerToNative.Compile(cliCtx.Context, compiler.CompilerInput{
				MainPkgPath:   compilePkg,
				OutputPath:    tempDir,
				EmitTraceFile: false,
				Mode:          compiler.ModeExecutable,
			}); err != nil {
				return err
			}

			// Compiler outputs binary as "output" (or "output.exe" on Windows)
			outputName := "output"
			if runtime.GOOS == "windows" {
				outputName += ".exe"
				binName += ".exe"
			}

			builtBinary := filepath.Join(tempDir, outputName)
			if _, err := os.Stat(builtBinary); err != nil {
				return fmt.Errorf("expected built binary at %s: %w", builtBinary, err)
			}

			// Install to GOBIN, GOPATH/bin, or ~/go/bin
			binDir, err := resolveBinDir()
			if err != nil {
				return err
			}

			if err := os.MkdirAll(binDir, 0o755); err != nil {
				return fmt.Errorf("create bin dir: %w", err)
			}

			targetPath := filepath.Join(binDir, binName)
			if err := os.Rename(builtBinary, targetPath); err != nil {
				return fmt.Errorf("install binary: %w", err)
			}

			fmt.Printf("installed %s to %s\n", binName, targetPath)
			return nil
		},
	}
}

// resolveBinDir determines the installation directory following Go conventions:
// GOBIN > GOPATH/bin > ~/go/bin
func resolveBinDir() (string, error) {
	if gobin := os.Getenv("GOBIN"); gobin != "" {
		return gobin, nil
	}

	if gopath := os.Getenv("GOPATH"); gopath != "" {
		return filepath.Join(gopath, "bin"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}

	return filepath.Join(home, "go", "bin"), nil
}

// resolveMainPkgPath resolves the actual package path, handling both module root
// and src subdirectory cases. If path is a file, returns it as-is. If it's a directory
// with .neva files, uses it. Otherwise checks for src subdirectory.
func resolveMainPkgPath(absPkg string) (string, error) {
	info, err := os.Stat(absPkg)
	if err != nil {
		return "", fmt.Errorf("stat package: %w", err)
	}

	if !info.IsDir() {
		return absPkg, nil
	}

	containsRootNeva, err := dirContainsNeva(absPkg)
	if err != nil {
		return "", fmt.Errorf("inspect package dir: %w", err)
	}
	if containsRootNeva {
		return absPkg, nil
	}

	// Check if src subdirectory contains .neva files
	srcDir := filepath.Join(absPkg, "src")
	srcInfo, err := os.Stat(srcDir)
	if err == nil && srcInfo.IsDir() {
		hasSrcNeva, err := dirContainsNeva(srcDir)
		if err != nil {
			return "", fmt.Errorf("inspect src dir: %w", err)
		}
		if hasSrcNeva {
			return srcDir, nil
		}
	}

	return absPkg, nil
}

func dirContainsNeva(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) == ".neva" {
			return true, nil
		}
	}

	return false, nil
}
