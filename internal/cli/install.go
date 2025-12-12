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

			if !pkgInfo.IsDir() {
				pkgBase = filepath.Dir(absPkg)
			}

			compilePkg, err := resolveMainPkgPath(absPkg)
			if err != nil {
				return err
			}

			binName := filepath.Base(pkgBase)
			if binName == string(filepath.Separator) || binName == "." || binName == "" {
				return fmt.Errorf("cannot determine binary name for %s", mainPkg)
			}

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

			outputName := "output"
			if runtime.GOOS == "windows" {
				outputName += ".exe"
				binName += ".exe"
			}

			builtBinary := filepath.Join(tempDir, outputName)
			if _, err := os.Stat(builtBinary); err != nil {
				return fmt.Errorf("expected built binary at %s: %w", builtBinary, err)
			}

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
