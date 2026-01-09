package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	cli "github.com/urfave/cli/v2"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/backend/golang"
	"github.com/nevalang/neva/internal/compiler/backend/golang/native"
	ir_backend "github.com/nevalang/neva/internal/compiler/backend/ir"
	"github.com/nevalang/neva/internal/compiler/desugarer"
)

func newRunCmd(
	workdir string,
	bldr builder.Builder,
	parser compiler.Parser,
	desugarer desugarer.Desugarer,
	analyzer compiler.Analyzer,
	irgen compiler.Irgen,
) *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Build and run neva program from source code",
		Args:  true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "emit-trace",
				Usage: "Write real-time trace to a file",
			},
			&cli.BoolFlag{
				Name:  "debug-runtime-validation",
				Usage: "Enable compiler runtime port validation (language developers only)",
			},
			&cli.BoolFlag{
				Name:  "emit-ir",
				Usage: "Emit intermediate representation before running",
			},
			&cli.StringFlag{
				Name:  "emit-ir-format",
				Usage: "Format for intermediate representation if IR is emitted",
			},
			&cli.BoolFlag{
				Name:  "watch",
				Usage: "Rebuild and rerun when source files change",
			},
		},
		ArgsUsage: "Provide path to main package",
		Action: func(cliCtx *cli.Context) error {
			emitIR := cliCtx.IsSet("emit-ir")
			isEmitIRFormatSet := cliCtx.IsSet("emit-ir-format")

			if !emitIR && isEmitIRFormatSet {
				return errors.New("emit-ir-format cannot be used without emit-ir")
			}

			var emitIRFormat ir_backend.Format
			if isEmitIRFormatSet {
				emitIRFormat = ir_backend.Format(cliCtx.String("emit-ir-format"))
			} else {
				emitIRFormat = ir_backend.FormatYAML
			}

			switch emitIRFormat {
			case ir_backend.FormatYAML, ir_backend.FormatJSON, ir_backend.FormatDOT, ir_backend.FormatMermaid, ir_backend.FormatThreeJS:
			default:
				return fmt.Errorf("unknown emit-ir-format: %s", emitIRFormat)
			}

			// we need to always set GOOS for compiler backend
			prevGOOS := os.Getenv("GOOS")
			if err := os.Setenv("GOOS", runtime.GOOS); err != nil {
				return fmt.Errorf("set GOOS: %w", err)
			}
			defer func() {
				if err := os.Setenv("GOOS", prevGOOS); err != nil {
					panic(err)
				}
			}()

			mainPkg, err := mainPkgPathFromArgs(cliCtx)
			if err != nil {
				return err
			}

			// Resolve mainPkg relative to workdir if it's not absolute
			if !filepath.IsAbs(mainPkg) {
				mainPkg = filepath.Join(workdir, mainPkg)
			}
			mainPkg = filepath.Clean(mainPkg)

			input := compiler.CompilerInput{
				MainPkgPath:   mainPkg,
				OutputPath:    workdir,
				EmitTraceFile: cliCtx.IsSet("emit-trace"),
				Mode:          compiler.ModeExecutable,
			}

			runOnce := func(ctx context.Context) error {
				tempExecDir, err := os.MkdirTemp("", "neva_run_")
				if err != nil {
					return fmt.Errorf("create temporary execution directory: %w", err)
				}
				defer os.RemoveAll(tempExecDir)

				if emitIR {
					irCompiler := compiler.New(
						bldr,
						parser,
						&desugarer,
						analyzer,
						irgen,
						ir_backend.NewBackend(emitIRFormat),
					)
					if _, err := irCompiler.Compile(ctx, input); err != nil {
						return fmt.Errorf("emit IR: %w", err)
					}
				}

				compilerToNative := compiler.New(
					bldr,
					parser,
					&desugarer,
					analyzer,
					irgen,
					native.NewBackend(
						golang.NewBackend("", cliCtx.Bool("debug-runtime-validation")),
					),
				)

				input.OutputPath = tempExecDir

				if _, err := compilerToNative.Compile(ctx, input); err != nil {
					return err
				}

				expectedOutputFileName := "output"
				if runtime.GOOS == "windows" { // assumption that on windows compiler generates .exe
					expectedOutputFileName += ".exe"
				}

				execPath := filepath.Join(tempExecDir, expectedOutputFileName)

				cmd := exec.CommandContext(ctx, execPath)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					return fmt.Errorf("failed to run generated executable: %w", err)
				}

				return nil
			}

			if cliCtx.Bool("watch") {
				moduleRoot, err := findModuleRoot(workdir, mainPkg)
				if err != nil {
					return err
				}

				return watchAndRun(cliCtx.Context, moduleRoot, runOnce)
			}

			return runOnce(cliCtx.Context)
		},
	}
}
