package cli

import (
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
)

func newRunCmd(
	workdir string,
	bldr builder.Builder,
	parser compiler.Parser,
	desugarer compiler.Desugarer,
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
				Name:  "emit-ir",
				Usage: "Emit intermediate representation before running",
			},
			&cli.StringFlag{
				Name:  "emit-ir-format",
				Usage: "Format for intermediate representation if IR is emitted",
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
			case ir_backend.FormatYAML, ir_backend.FormatJSON:
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

			input := compiler.CompilerInput{
				MainPkgPath:   mainPkg,
				OutputPath:    workdir,
				EmitTraceFile: cliCtx.IsSet("emit-trace"),
			}

			compilerToNative := compiler.New(
				bldr,
				parser,
				desugarer,
				analyzer,
				irgen,
				native.NewBackend(
					golang.NewBackend(),
				),
			)

			out, err := compilerToNative.Compile(cliCtx.Context, input)
			if err != nil {
				return err
			}

			irBackend := ir_backend.NewBackend(emitIRFormat)
			// TODO refactor - trace is only used by golang and golang/native backends
			// it should not be part of the compiler.Backend interface.
			if err := irBackend.Emit(workdir, out.MiddleEnd.IR, false); err != nil {
				return err
			}

			expectedOutputFileName := "output"
			if runtime.GOOS == "windows" { // assumption that on windows compiler generates .exe
				expectedOutputFileName += ".exe"
			}

			execPath := filepath.Join(workdir, expectedOutputFileName)

			defer func() {
				if err := os.Remove(execPath); err != nil {
					fmt.Println("failed to remove output file:", err)
				}
			}()

			cmd := exec.CommandContext(cliCtx.Context, execPath)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to run generated executable: %w", err)
			}

			return nil
		},
	}
}
