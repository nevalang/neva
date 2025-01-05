package cli

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler"

	cli "github.com/urfave/cli/v2"
)

func newBuildCmd(
	workdir string,
	compilerToGo compiler.Compiler,
	compilerToNative compiler.Compiler,
	compilerToWASM compiler.Compiler,
	compilerToJSON compiler.Compiler,
	compilerToDOT compiler.Compiler,
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
				Usage: "Target platform for build (options: go, wasm, native, json, dot)",
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
			mainPkg, err := mainPkgPathFromArgs(cliCtx)
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
			} else {
				target = "native"
			}

			var isTraceEnabled bool
			if cliCtx.IsSet("trace") {
				isTraceEnabled = true
			}

			compilerInput := compiler.CompilerInput{
				Main:   mainPkg,
				Output: output,
				Trace:  isTraceEnabled,
			}

			var compilerToUse compiler.Compiler
			switch target {
			case "go":
				compilerToUse = compilerToGo
			case "wasm":
				compilerToUse = compilerToWASM
			case "json":
				compilerToUse = compilerToJSON
			case "dot":
				compilerToUse = compilerToDOT
			case "native":
				compilerToUse = compilerToNative
			default:
				return fmt.Errorf("Unknown target %s", target)
			}

			return compilerToUse.Compile(cliCtx.Context, compilerInput)
		},
	}
}
