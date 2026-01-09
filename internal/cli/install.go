package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	cli "github.com/urfave/cli/v2"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/backend/golang"
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

			// Generate Go source to temporary directory
			tempDir, err := os.MkdirTemp("", "neva-install-*")
			if err != nil {
				return fmt.Errorf("create temp dir: %w", err)
			}
			defer os.RemoveAll(tempDir)

			// Use golang backend to generate Go source (not native, which builds)
			compilerToGo := compiler.New(
				bldr,
				parser,
				&desugarer,
				analyzer,
				irgen,
				golang.NewBackend("", false),
			)

			if _, err := compilerToGo.Compile(cliCtx.Context, compiler.CompilerInput{
				MainPkgPath:   compilePkg,
				OutputPath:    tempDir,
				EmitTraceFile: false,
				Mode:          compiler.ModeExecutable,
			}); err != nil {
				return err
			}

			// Determine installation directory (GOBIN, GOPATH/bin, or ~/go/bin)
			binDir, err := resolveBinDir()
			if err != nil {
				return err
			}

			// Create bin directory if it doesn't exist
			if err := os.MkdirAll(binDir, 0o755); err != nil {
				return fmt.Errorf("create bin dir: %w", err)
			}

			if runtime.GOOS == "windows" {
				binName += ".exe"
			}

			targetPath := filepath.Join(binDir, binName)

			// Use go build to build directly to the target location
			// This leverages Go's build system while giving us control over binary name
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("get working directory: %w", err)
			}

			// Change to temporary directory to build the binary
			if err := os.Chdir(tempDir); err != nil {
				return fmt.Errorf("change directory to temp: %w", err)
			}
			defer func() {
				if err := os.Chdir(wd); err != nil {
					panic(err)
				}
			}()

			cmd := exec.Command("go", "build", "-ldflags", "-s -w", "-o", targetPath, ".")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("go build: %w", err)
			}

			if _, err := os.Stat(targetPath); err != nil {
				return fmt.Errorf("binary not found at %s after build: %w", targetPath, err)
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
