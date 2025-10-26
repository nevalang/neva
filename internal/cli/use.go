package cli

import (
	"fmt"

	cli "github.com/urfave/cli/v2"

	"github.com/nevalang/neva/internal/versionmanager"
	"github.com/nevalang/neva/pkg"
)

func newUseCmd() *cli.Command {
	return &cli.Command{
		Name:        "use",
		Usage:       "Install (if needed) and activate a specific neva release",
		ArgsUsage:   "<version>|latest",
		Description: "Downloads the requested release (if missing), marks it as active, and proxies future invocations of this bundled CLI to it.",
		Action: func(cCtx *cli.Context) error {
			if cCtx.Args().Len() != 1 {
				return fmt.Errorf("expected 1 argument, got %d", cCtx.Args().Len())
			}

			manager, err := versionmanager.NewManager()
			if err != nil {
				return err
			}

			currentTag, err := versionmanager.Normalize(pkg.Version)
			if err != nil {
				return err
			}

			normalized, installed, err := manager.Use(cCtx.Context, cCtx.Args().First(), pkg.Version)
			if err != nil {
				return err
			}

			if normalized == currentTag {
				fmt.Printf("Using bundled neva %s\n", normalized)
				return nil
			}

			if installed {
				fmt.Printf("Installed neva %s\n", normalized)
			} else {
				fmt.Printf("neva %s was already installed\n", normalized)
			}

			fmt.Printf("Now using neva %s\n", normalized)
			fmt.Println("Tip: subsequent 'neva' invocations will automatically proxy to the selected release.")

			return nil
		},
	}
}
