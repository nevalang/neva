package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/urfave/cli/v2"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/backend/golang"
	"github.com/nevalang/neva/internal/compiler/backend/golang/native"
	"github.com/nevalang/neva/internal/compiler/backend/golang/wasm"
	"github.com/nevalang/neva/internal/compiler/backend/ir"
	"github.com/nevalang/neva/internal/compiler/desugarer"
	"github.com/nevalang/neva/internal/compiler/irgen"
	"github.com/nevalang/neva/internal/compiler/parser"
)

//nolint:gocyclo // CLI flag setup is dense; refactor later without behavior changes.
func newBuildCmd(
	workdir string,
	bldr builder.Builder,
	parser parser.Parser,
	desugarer desugarer.Desugarer,
	analyzer analyzer.Analyzer,
	irgen irgen.Generator,
) *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Build neva program",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "output",
				Usage: "Output directory",
			},
			&cli.BoolFlag{
				Name:  "emit-trace",
				Usage: "Emit trace file",
			},
			&cli.BoolFlag{
				Name:  "debug-runtime-validation",
				Usage: "Enable compiler runtime port validation (language developers only)",
			},
			&cli.StringFlag{
				Name:  "target",
				Usage: "Target platform (go, wasm, native, ir)",
				Value: "native",
			},
			&cli.StringFlag{
				Name:  "target-os",
				Usage: "Target OS (only for native target)",
			},
			&cli.StringFlag{
				Name:  "target-arch",
				Usage: "Target Architecture (only for native target)",
			},
			&cli.StringFlag{
				Name:  "target-go-mode",
				Usage: "Go target mode (executable, pkg)",
				Value: "executable",
			},
			&cli.StringFlag{
				Name:  "target-ir-format",
				Usage: "IR target format (json, yaml, dot, mermaid, ascii, threejs)",
				Value: "yaml",
			},
			&cli.StringFlag{
				Name:  "target-go-runtime-path",
				Usage: "Go runtime import path (only for go target)",
			},
		},
		ArgsUsage: "Provide path to main package",
		Action: func(cliCtx *cli.Context) error {
			var target string
			if cliCtx.IsSet("target") {
				target = cliCtx.String("target")
			} else {
				target = "native"
			}

			switch target {
			case "go", "wasm", "ir", "native":
			default:
				return fmt.Errorf("unknown target %s", target)
			}

			targetOS := cliCtx.String("target-os")
			if targetOS != "" && target != "native" {
				return fmt.Errorf("target-os and target-arch are only supported when target is native")
			}

			targetArch := cliCtx.String("target-arch")
			if targetArch != "" && target != "native" {
				return fmt.Errorf("target-arch is only supported when target is native")
			}

			if (targetOS != "" && targetArch == "") || (targetOS == "" && targetArch != "") {
				return fmt.Errorf("target-os and target-arch must be set together")
			}

			isIRTargetFormatSet := cliCtx.IsSet("target-ir-format")
			if isIRTargetFormatSet && target != "ir" {
				return errors.New("target-ir-format cannot be used when target is not ir")
			}

			var irTargetFormat ir.Format
			if isIRTargetFormatSet {
				irTargetFormat = ir.Format(cliCtx.String("target-ir-format"))
			} else {
				irTargetFormat = ir.FormatYAML
			}

			switch irTargetFormat {
			case ir.FormatYAML, ir.FormatJSON, ir.FormatDOT, ir.FormatMermaid, ir.FormatASCII, ir.FormatThreeJS:
			default:
				return fmt.Errorf("unknown target-ir-format: %s", irTargetFormat)
			}

			mainPkgPath, err := mainPkgPathFromArgs(cliCtx, workdir)
			if err != nil {
				return err
			}

			outputDirPath := workdir
			if cliCtx.IsSet("output") {
				outputDirPath = cliCtx.String("output")
			}

			// we're going to change GOOS and GOARCH, so we need to restore them after compilation
			prevGOOS := os.Getenv("GOOS")
			prevGOARCH := os.Getenv("GOARCH")
			// if target-os and target-arch are not set, use the current platform
			if targetOS == "" {
				targetOS = runtime.GOOS
				targetArch = runtime.GOARCH
			}
			// compiler backend (native one) depends on GOOS and GOARCH, so we always must set them
			if err := os.Setenv("GOOS", targetOS); err != nil {
				return fmt.Errorf("set GOOS: %w", err)
			}
			if err := os.Setenv("GOARCH", targetArch); err != nil {
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

			var compilerToUse compiler.Compiler
			var compilerMode compiler.Mode

			externalRuntimePath := cliCtx.String("target-go-runtime-path")
			if externalRuntimePath != "" && !filepath.IsAbs(externalRuntimePath) {
				externalRuntimePath = filepath.Join(workdir, externalRuntimePath)
			}

			golangBackend := golang.NewBackend(
				externalRuntimePath,
				cliCtx.Bool("debug-runtime-validation"),
			)

			switch target {
			case "go":
				goMode := cliCtx.String("target-go-mode")
				switch goMode {
				case "executable", "":
					compilerMode = compiler.ModeExecutable
				case "pkg":
					compilerMode = compiler.ModeLibrary
				default:
					return fmt.Errorf("unknown target-go-mode: %s", goMode)
				}
				compilerToUse = compiler.New(
					bldr,
					parser,
					&desugarer,
					analyzer,
					irgen,
					golangBackend,
				)
			case "wasm":
				compilerMode = compiler.ModeExecutable
				compilerToUse = compiler.New(
					bldr,
					parser,
					&desugarer,
					analyzer,
					irgen,
					wasm.NewBackend(golangBackend),
				)
			case "ir":
				compilerMode = compiler.ModeExecutable
				compilerToUse = compiler.New(
					bldr,
					parser,
					&desugarer,
					analyzer,
					irgen,
					ir.NewBackend(irTargetFormat),
				)
			case "native":
				compilerMode = compiler.ModeExecutable
				compilerToUse = compiler.New(
					bldr,
					parser,
					&desugarer,
					analyzer,
					irgen,
					native.NewBackend(golangBackend),
				)
			}

			if _, err := compilerToUse.Compile(cliCtx.Context, compiler.CompilerInput{
				MainPkgPath:   mainPkgPath,
				OutputPath:    outputDirPath,
				EmitTraceFile: cliCtx.IsSet("emit-trace"),
				Mode:          compilerMode,
			}); err != nil {
				return fmt.Errorf("failed to compile: %w", err)
			}

			return nil
		},
	}
}

func mainPkgPathFromArgs(cliCtx *cli.Context, workdir string) (string, error) {
	if cliCtx.NArg() == 0 {
		return "", errors.New("path to main package is required")
	}
	raw := cliCtx.Args().First()

	// Resolve path relative to workdir if not absolute
	abs := raw
	if !filepath.IsAbs(abs) {
		abs = filepath.Join(workdir, raw)
	}
	abs = filepath.Clean(abs)

	// Check if path exists and is a .neva file - if so, use its directory
	info, err := os.Stat(abs)
	if err == nil && !info.IsDir() {
		if filepath.Ext(abs) == ".neva" {
			abs = filepath.Dir(abs)
		}
	}

	return abs, nil
}
