package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	cli "github.com/urfave/cli/v2"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/interpreter"
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
	var (
		target           string
		debug            bool
		debugLogFilePath string
	)

	return &cli.App{
		Name:  "neva",
		Usage: "Flow-based programming language",
		Commands: []*cli.Command{
			{
				Name:  "version",
				Usage: "Get current Nevalang version",
				Action: func(_ *cli.Context) error {
					fmt.Println(pkg.Version)
					return nil
				},
			},
			{
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
			},
			{
				Name:  "new",
				Usage: "Create new Nevalang project",
				Args:  true,
				Action: func(cCtx *cli.Context) error {
					if path := cCtx.Args().First(); path != "" {
						if err := os.Mkdir(path, 0755); err != nil {
							return err
						}
						return createNevaMod(path)
					}
					return createNevaMod(workdir)
				},
			},
			{
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
			},
			{
				Name:      "run",
				Usage:     "Run neva program from source code in interpreter mode",
				Args:      true,
				ArgsUsage: "Provide path to the executable package",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "debug",
						Usage:       "Show message events in stdout",
						Destination: &debug,
					},
					&cli.StringFlag{
						Name:        "debugLogFilePath",
						Usage:       "File path to write debug log (only available if -debug is passed)",
						Destination: &debugLogFilePath,
					},
				},
				Action: func(cCtx *cli.Context) error {
					if !debug && debugLogFilePath != "" {
						return fmt.Errorf("debugFile can only be used with -debug flag")
					}

					dirFromArg, err := getMainPkgFromArgs(cCtx)
					if err != nil {
						return err
					}

					return interpreter.New(bldr, goc).
						Interpret(context.Background(), dirFromArg, debug, debugLogFilePath)
				},
			},
			{
				Name:  "build",
				Usage: "Build executable binary from neva program source code",
				Args:  true,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "target",
						Required:    false,
						Usage:       "Emit Go or WASM instead of machine code",
						Destination: &target,
						Action: func(ctx *cli.Context, s string) error {
							switch s {
							case "go", "wasm", "native", "json", "dot":
							default:
								return fmt.Errorf("Unknown target %s", s)
							}
							return nil
						},
					},
				},
				ArgsUsage: "Provide path to the executable package",
				Action: func(cCtx *cli.Context) error {
					dirFromArg, err := getMainPkgFromArgs(cCtx)
					if err != nil {
						return err
					}
					switch target {
					case "go":
						return goc.Compile(dirFromArg, workdir, debug)
					case "wasm":
						return wasmc.Compile(dirFromArg, workdir, debug)
					case "json":
						return jsonc.Compile(dirFromArg, workdir, debug)
					case "dot":
						return dotc.Compile(dirFromArg, workdir, debug)
					default:
						return nativec.Compile(dirFromArg, workdir, debug)
					}
				},
			},
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

func createNevaMod(path string) error {
	// Create neva.yml file
	nevaYmlContent := fmt.Sprintf("neva: %s", pkg.Version)
	if err := os.WriteFile(
		filepath.Join(path, "neva.yml"),
		[]byte(nevaYmlContent),
		0644,
	); err != nil {
		return err
	}

	// Create src sub-directory
	srcPath := filepath.Join(path, "src")
	if err := os.Mkdir(srcPath, 0755); err != nil {
		return err
	}

	// Create main.neva file
	mainNevaContent := `flow Main(start) (stop) {
	nodes { Println }
	:start -> ('Hello, World!' -> println -> :stop)
}`

	if err := os.WriteFile(
		filepath.Join(srcPath, "main.neva"),
		[]byte(mainNevaContent),
		0644,
	); err != nil {
		return err
	}

	return nil
}
