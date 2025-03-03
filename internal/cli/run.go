package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/nevalang/neva/internal/compiler"

	cli "github.com/urfave/cli/v2"
)

func newRunCmd(workdir string, nativec compiler.Compiler) *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Build and run neva program from source code",
		Args:  true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "trace",
				Usage: "Write trace information to file",
			},
			&cli.BoolFlag{
				Name:  "emit-ir",
				Usage: "Emit intermediate representation to ir.yml file",
			},
		},
		ArgsUsage: "Provide path to main package",
		Action: func(cliCtx *cli.Context) error {
			mainPkg, err := mainPkgPathFromArgs(cliCtx)
			if err != nil {
				return err
			}

			trace := cliCtx.IsSet("trace")
			emitIR := cliCtx.IsSet("emit-ir")

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

			expectedOutputFileName := "output"
			if runtime.GOOS == "windows" { // assumption that on windows compiler generates .exe
				expectedOutputFileName += ".exe"
			}

			input := compiler.CompilerInput{
				Main:   mainPkg,
				Output: workdir,
				Trace:  trace,
				EmitIR: emitIR,
			}

			if err := nativec.Compile(cliCtx.Context, input); err != nil {
				return err
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
