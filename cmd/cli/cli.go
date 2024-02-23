package main

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/interpreter"
	"github.com/nevalang/neva/pkg"
	cli "github.com/urfave/cli/v2"
)

func newCliApp( //nolint:funlen
	wd string,
	goComp compiler.Compiler,
	nativeComp compiler.Compiler,
	wasmComp compiler.Compiler,
	intr interpreter.Interpreter,
) *cli.App {
	var target string

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
				Name:      "run",
				Usage:     "Run neva program from source code in interpreter mode",
				Args:      true,
				ArgsUsage: "Provide path to the executable package",
				Action: func(cCtx *cli.Context) error {
					dirFromArg, err := getMainPkgFromArgs(cCtx)
					if err != nil {
						return err
					}
					if err := intr.Interpret(
						context.Background(),
						wd,
						dirFromArg,
					); err != nil {
						return err
					}
					return nil
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
							case "go", "wasm", "native":
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
						return goComp.Compile(
							wd, dirFromArg, wd,
						)
					case "wasm":
						return wasmComp.Compile(
							wd, dirFromArg, wd,
						)
					default:
						return nativeComp.Compile(
							wd, dirFromArg, wd,
						)
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
