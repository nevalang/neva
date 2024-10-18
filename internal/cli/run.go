package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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
		},
		ArgsUsage: "Provide path to main package",
		Action: func(cliCtx *cli.Context) error {
			mainPkg, err := getMainPkgFromArgs(cliCtx)
			if err != nil {
				return err
			}

			output := workdir
			if cliCtx.IsSet("output") {
				output = cliCtx.String("output")
			}

			var trace bool
			if cliCtx.IsSet("trace") {
				trace = true
			}

			input := compiler.CompilerInput{
				Main:   mainPkg,
				Output: output,
				Trace:  trace,
			}

			if err := nativec.Compile(cliCtx.Context, input); err != nil {
				return err
			}

			defer func() {
				if err := os.Remove(filepath.Join(workdir, "output")); err != nil {
					fmt.Println("failed to remove output file:", err)
				}
			}()

			pathToExec := filepath.Join(workdir, "output")

			cmd := exec.CommandContext(cliCtx.Context, pathToExec)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			return cmd.Run()
		},
	}
}
