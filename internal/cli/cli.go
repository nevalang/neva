package cli

import (
	"fmt"
	"os/exec"
	"runtime"

	cli "github.com/urfave/cli/v2"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/desugarer"
	"github.com/nevalang/neva/internal/compiler/irgen"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/pkg"
)

func NewApp(
	workdir string,
	bldr builder.Builder,
	prsr parser.Parser,
	desugarer desugarer.Desugarer,
	analyzer analyzer.Analyzer,
	irgen irgen.Generator,
) *cli.App {
	return &cli.App{
		Name:  "neva",
		Usage: "Dataflow programming language with static types and implicit parallelism",
		Commands: []*cli.Command{
			newVersionCmd(),
			newFmtCmd(),
			newUpgradeCmd(),
			newToolCmd(),
			newNewCmd(),
			newGetCmd(workdir, bldr),
			newInstallCmd(workdir, bldr, prsr, desugarer, analyzer, irgen),
			newRunCmd(workdir, bldr, prsr, desugarer, analyzer, irgen),
			newBuildCmd(workdir, bldr, prsr, desugarer, analyzer, irgen),
			newOSArchCmd(),
			newDocCmd(),
		},
	}
}

func newVersionCmd() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Get current Nevalang version",
		Action: func(_ *cli.Context) error {
			fmt.Println(pkg.Version)
			return nil
		},
	}
}

func newUpgradeCmd() *cli.Command {
	return &cli.Command{
		Name:  "upgrade",
		Usage: "Upgrade to newest Nevalang version",
		Action: func(cliCtx *cli.Context) error {
			command, args := upgradeCommand(runtime.GOOS)
			err := exec.CommandContext(cliCtx.Context, command, args...).Run()
			if err != nil {
				fmt.Println("Upgrading Nevalang failed: " + err.Error())
			} else {
				fmt.Println("Upgrading Nevalang completed. Restart your terminal or VS Code, then run `neva version` to verify the installed version.")
			}
			return nil
		},
	}
}

func upgradeCommand(goos string) (string, []string) {
	const unixInstaller = "https://raw.githubusercontent.com/nevalang/neva/main/scripts/install.sh"
	if goos == "windows" {
		const windowsInstaller = "https://raw.githubusercontent.com/nevalang/neva/main/scripts/install.bat"
		return "powershell", []string{
			"-NoProfile",
			"-Command",
			"$script = Join-Path $env:TEMP 'neva-install.bat'; " +
				"Invoke-WebRequest '" + windowsInstaller + "' -OutFile $script; " +
				"cmd /c $script",
		}
	}

	return "sh", []string{"-c", "curl -fsSL " + unixInstaller + " | sh"}
}

func newGetCmd(workdir string, bldr builder.Builder) *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "Add dependency to current module",
		Args:      true,
		ArgsUsage: "Provide path to the module",
		Action: func(cliCtx *cli.Context) error {
			if cliCtx.Args().Len() != 2 {
				return fmt.Errorf(
					"expected 2 arguments, got %d",
					cliCtx.Args().Len(),
				)
			}

			path := cliCtx.Args().Get(0)
			version := cliCtx.Args().Get(1)

			installedPath, err := bldr.Get(workdir, path, version)
			if err != nil {
				return fmt.Errorf("failed to get dependency: %w", err)
			}

			fmt.Printf(
				"%s installed to %s\n", cliCtx.Args().Get(0),
				installedPath,
			)

			return nil
		},
	}
}
