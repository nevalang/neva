package cli

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	cli "github.com/urfave/cli/v2"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/pkg"
)

func NewApp(
	workdir string,
	bldr builder.Builder,
	goc compiler.Compiler,
	nativec compiler.Compiler,
	wasmc compiler.Compiler,
	jsonc compiler.Compiler,
	dotc compiler.Compiler,
) *cli.App {
	return &cli.App{
		Name:  "neva",
		Usage: "Dataflow programming language with static types and implicit parallelism",
		Commands: []*cli.Command{
			versionCmd,
			upgradeCmd,
			newNewCmd(workdir),
			newGetCmd(workdir, bldr),
			newRunCmd(workdir, nativec),
			newBuildCmd(workdir, goc, nativec, wasmc, jsonc, dotc),
		},
	}
}

var versionCmd = &cli.Command{
	Name:  "version",
	Usage: "Get current Nevalang version",
	Action: func(_ *cli.Context) error {
		fmt.Println(pkg.Version)
		return nil
	},
}

var upgradeCmd = &cli.Command{
	Name:  "upgrade",
	Usage: "Upgrade to newest Nevalang version",
	Action: func(_ *cli.Context) error {
		cmd := exec.Command("curl -sSL https://raw.githubusercontent.com/nevalang/neva/main/scripts/install.sh | bash")
		err := cmd.Run()
		if err != nil {
			fmt.Println("Upgrading Nevalang failed :" + err.Error())
		} else {
			fmt.Println("Upgrading Nevalang completed. Upgraded to version: " + pkg.Version)
		}
		return nil
	},
}

func newGetCmd(workdir string, bldr builder.Builder) *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "Add dependency to current module",
		Args:      true,
		ArgsUsage: "Provide path to the module",
		Action: func(cCtx *cli.Context) error {
			installedPath, err := bldr.Get(
				workdir,
				cCtx.Args().Get(0),
				cCtx.Args().Get(1),
			)
			if err != nil {
				return err
			}
			fmt.Printf(
				"%s installed to %s\n", cCtx.Args().Get(0),
				installedPath,
			)
			return nil
		},
	}
}

func getMainPkgFromArgs(cCtx *cli.Context) (string, error) {
	firstArg := cCtx.Args().First()
	dirFromArg := strings.TrimSuffix(firstArg, "/main.neva")
	if filepath.Ext(dirFromArg) != "" {
		return "", errors.New(
			"Use path to directory with executable package, relative to module root",
		)
	}
	return dirFromArg, nil
}
