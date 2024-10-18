package cli

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"

	cli "github.com/urfave/cli/v2"
)

func newBuildCmd(
	workdir string,
	goc compiler.Compiler,
	nativec compiler.Compiler,
	wasmc compiler.Compiler,
	jsonc compiler.Compiler,
	dotc compiler.Compiler,
) *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Build neva program from source code",
		Args:  true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "output",
				Usage: "Where to put output file(s)",
			},
			&cli.BoolFlag{
				Name:  "trace",
				Usage: "Write trace information to file",
			},
			&cli.StringFlag{
				Name:  "target",
				Usage: "Target platform for build",
				Action: func(ctx *cli.Context, s string) error {
					switch s {
					case "go", "wasm", "native", "json", "dot":
						return nil
					}
					return fmt.Errorf("Unknown target %s", s)
				},
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

			var target string
			if cliCtx.IsSet("target") {
				target = cliCtx.String("target")
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

			var c compiler.Compiler
			switch target {
			case "go":
				c = goc
			case "wasm":
				c = wasmc
			case "json":
				c = jsonc
			case "dot":
				c = dotc
			default:
				c = nativec
			}

			return c.Compile(cliCtx.Context, input)
		},
	}
}
