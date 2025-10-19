package cli

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/backend/dot"
	"github.com/nevalang/neva/internal/compiler/backend/golang"
	"github.com/nevalang/neva/internal/compiler/backend/golang/native"
	"github.com/nevalang/neva/internal/compiler/backend/golang/wasm"
	"github.com/nevalang/neva/internal/compiler/backend/ir"
	ir_backend "github.com/nevalang/neva/internal/compiler/backend/ir"

	cli "github.com/urfave/cli/v2"
)

func newBuildCmd(
	workdir string,
	bldr builder.Builder,
	parser compiler.Parser,
	desugarer compiler.Desugarer,
	analyzer compiler.Analyzer,
	irgen compiler.Irgen,
) *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Generate target platform code from neva program",
		Args:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "output",
				Usage: "Where to put output file(s)",
			},
			&cli.BoolFlag{
				Name:  "emit-trace",
				Usage: "Emit trace.log file when running the program",
			},
			&cli.StringFlag{
				Name:  "target",
				Usage: "Target platform for build (options: go, wasm, native, json, dot). For 'native' target, 'target-os' and 'target-arch' flags can be used, but if used, they must be used together.",
				Action: func(ctx *cli.Context, s string) error {
					switch s {
					case "go", "wasm", "native", "ir", "dot":
						return nil
					}
					return fmt.Errorf("Unknown target %s", s)
				},
			},
			&cli.StringFlag{
				Name:  "target-os",
				Usage: "Target operating system for native build. See 'neva osarch' for supported combinations. Only supported for native target. Not needed if building for the current platform. Must be combined properly with 'target-arch'.",
			},
			&cli.StringFlag{
				Name:  "target-arch",
				Usage: "Target architecture for native build. See 'neva osarch' for supported combinations. Only supported for native target. Not needed if building for the current platform. Must be combined properly with 'target-os'.",
			},
			&cli.StringFlag{
				Name:  "target-ir-format",
				Usage: "Format for ir file - yaml or json",
			},
			&cli.StringFlag{
				Name:  "target-go-mode",
				Usage: "Go backend mode when target=go (executable|pkg)",
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
			case "go", "wasm", "ir", "dot", "native":
			default:
				return fmt.Errorf("Unknown target %s", target)
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

			var irTargetFormat ir_backend.Format
			if isIRTargetFormatSet {
				irTargetFormat = ir_backend.Format(cliCtx.String("target-ir-format"))
			} else {
				irTargetFormat = ir_backend.FormatYAML
			}

			switch irTargetFormat {
			case ir_backend.FormatYAML, ir_backend.FormatJSON:
			default:
				return fmt.Errorf("unknown target-ir-format: %s", irTargetFormat)
			}

			mainPkgPath, err := mainPkgPathFromArgs(cliCtx)
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
			switch target {
			case "go":
				switch goMode := cliCtx.String("target-go-mode"); goMode {
				case "executable", "", "pkg":
					break
				default:
					return fmt.Errorf("unknown target-go-mode: %s", goMode)
				}
				compilerToUse = compiler.New(
					bldr,
					parser,
					desugarer,
					analyzer,
					irgen,
					golang.NewBackend(
						golang.Mode(
							cliCtx.String("target-go-mode"),
						),
					),
				)
			case "wasm":
				compilerToUse = compiler.New(
					bldr,
					parser,
					desugarer,
					analyzer,
					irgen,
					wasm.NewBackend(
						golang.NewBackend(golang.ModeExecutable),
					),
				)
			case "ir":
				compilerToUse = compiler.New(
					bldr,
					parser,
					desugarer,
					analyzer,
					irgen,
					ir.NewBackend(irTargetFormat),
				)
			case "dot":
				compilerToUse = compiler.New(
					bldr,
					parser,
					desugarer,
					analyzer,
					irgen,
					dot.NewBackend(),
				)
			case "native":
				compilerToUse = compiler.New(
					bldr,
					parser,
					desugarer,
					analyzer,
					irgen,
					native.NewBackend(golangBackend),
				)
			}

			if _, err := compilerToUse.Compile(cliCtx.Context, compiler.CompilerInput{
				MainPkgPath:   mainPkgPath,
				OutputPath:    outputDirPath,
				EmitTraceFile: cliCtx.IsSet("emit-trace"),
			}); err != nil {
				return fmt.Errorf("failed to compile: %w", err)
			}

			return nil
		},
	}
}
