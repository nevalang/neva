package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	cli "github.com/urfave/cli/v2"
)

func newToolCmd() *cli.Command {
	return &cli.Command{
		Name:      "tool",
		Usage:     "Run an installed Neva developer tool",
		ArgsUsage: "<name> [arguments...]",
		Action: func(cliCtx *cli.Context) error {
			if cliCtx.Args().Len() == 0 {
				return errors.New("usage: neva tool <name> [arguments...]")
			}

			toolName := cliCtx.Args().First()
			binaryName := "neva-" + toolName
			if runtime.GOOS == "windows" {
				binaryName += ".exe"
			}
			binaryPath, err := exec.LookPath(binaryName)
			if err != nil {
				return fmt.Errorf("find %s in PATH: %w", binaryName, err)
			}

			command := exec.CommandContext(cliCtx.Context, binaryPath, cliCtx.Args().Tail()...)
			command.Stdin = os.Stdin
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			if err := command.Run(); err != nil {
				return fmt.Errorf("run %s: %w", binaryName, err)
			}
			return nil
		},
	}
}
