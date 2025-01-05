package cli

import (
	"fmt"
	"os/exec"
	"strings"

	cli "github.com/urfave/cli/v2"
)

func newOSArchCmd() *cli.Command {
	return &cli.Command{
		Name:  "osarch",
		Usage: "List supported OS/architecture combinations for native target",
		Action: func(cliCtx *cli.Context) error {
			cmd := exec.Command("go", "tool", "dist", "list")
			output, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("failed to execute go tool dist list: %w", err)
			}

			fmt.Println("Supported OS/architecture combinations for native target:")
			fmt.Println("(use these values for --target-os and --target-arch flags when cross-compiling)")
			fmt.Println()

			platforms := strings.Split(strings.TrimSpace(string(output)), "\n")
			for _, platform := range platforms {
				fmt.Println(platform)
			}

			return nil
		},
	}
}
