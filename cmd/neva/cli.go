package main

import (
	"context"
	"errors"
	"path/filepath"
	"strings"

	"github.com/nevalang/neva/internal/interpreter"
	cli "github.com/urfave/cli/v2"
)

func newApp(intr interpreter.Interpreter, wd string) *cli.App {
	return &cli.App{
		Name:  "neva",
		Usage: "Flow-based programming language",
		Commands: []*cli.Command{
			{
				Name:      "run",
				Usage:     "Run neva program from source code in interpreter mode",
				Args:      true,
				ArgsUsage: "Provide path to the executable package",
				Action: func(cCtx *cli.Context) error {
					path := strings.TrimSuffix(cCtx.Args().First(), "/main.neva")
					if filepath.Ext(path) != "" {
						return errors.New(
							"Use path to directory with executable package, relative to module root",
						)
					}
					if err := intr.Interpret(
						context.Background(),
						wd,
						path,
					); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:      "build",
				Usage:     "Build executable binary from neva program source code",
				Args:      true,
				ArgsUsage: "Provide path to the executable package",
				Action: func(cCtx *cli.Context) error {
					panic("not implemented")
				},
			},
		},
	}
}
