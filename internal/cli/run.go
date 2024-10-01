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

			if err := nativec.Compile(mainPkg, output, trace); err != nil {
				return fmt.Errorf("build failed: %w", err)
			}

			cmd := exec.Command(filepath.Join(workdir, "output"))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("run failed: %w", err)
			}

			// TODO cleanup

			return nil
		},
	}
}
