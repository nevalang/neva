package cli

import (
	"fmt"
	"os"

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
		Usage: "Generate target platform code from neva program",
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
				Usage: "Target platform for build (options: go, wasm, native, json, dot). For 'native' target, 'target-os' and 'target-arch' flags can be used, but if used, they must be used together.",
				Action: func(ctx *cli.Context, s string) error {
					switch s {
					case "go", "wasm", "native", "json", "dot":
						return nil
					}
					return fmt.Errorf("Unknown target %s", s)
				},
			},
			&cli.StringFlag{
				Name:  "target-os",
				Usage: "Target operating system for native build. See 'neva osarch' for supported combinations. Only supported for native target. Not needed if building for the current platform. Must be combined properly with 'target-arch'.",
			},
			&cli.StringFlag{
				Name:  "target-arch",
				Usage: "Target architecture for native build. See 'neva osarch' for supported combinations. Only supported for native target. Not needed if building for the current platform. Must be combined properly with 'target-os'.",
			},
		},
		ArgsUsage: "Provide path to main package",
		Action: func(cliCtx *cli.Context) error {
			var target string
			if cliCtx.IsSet("target") {
				target = cliCtx.String("target")
			} else {
				target = "native"
			}

			switch target {
			case "go", "wasm", "json", "dot", "native":
			default:
				return fmt.Errorf("Unknown target %s", target)
			}

			var targetOS, targetArch string
			if cliCtx.IsSet("target-os") {
				targetOS = cliCtx.String("target-os")
			}
			if cliCtx.IsSet("target-arch") {
				targetArch = cliCtx.String("target-arch")
			}
			if (targetOS != "" && targetArch == "") || (targetOS == "" && targetArch != "") {
				return fmt.Errorf("target-os and target-arch must be set together")
			}
			if target != "native" && targetOS != "" {
				return fmt.Errorf("target-os and target-arch are only supported when target is native")
			}

			mainPkg, err := mainPkgPathFromArgs(cliCtx)
			if err != nil {
				return err
			}

			outputDirPath := workdir
			if cliCtx.IsSet("output") {
				outputDirPath = cliCtx.String("output")
			}

			var isTraceEnabled bool
			if cliCtx.IsSet("trace") {
				isTraceEnabled = true
			}

			compilerInput := compiler.CompilerInput{
				Main:   mainPkg,
				Output: outputDirPath,
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
			}

			if targetOS != "" {
				prevGOOS := os.Getenv("GOOS")
				prevGOARCH := os.Getenv("GOARCH")

				if err := os.Setenv("GOOS", targetOS); err != nil {
					return fmt.Errorf("set GOOS: %w", err)
				}
				if err := os.Setenv("GOARCH", targetArch); err != nil {
					return fmt.Errorf("set GOARCH: %w", err)
				}

				defer func() {
					if err := os.Setenv("GOOS", prevGOOS); err != nil {
						panic(err)
					}
					if err := os.Setenv("GOARCH", prevGOARCH); err != nil {
						panic(err)
					}
				}()
			}

			return compilerToUse.Compile(cliCtx.Context, compilerInput)
		},
	}
}
