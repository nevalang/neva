package main

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	cli "github.com/urfave/cli/v2"
)

func main() {
	intr := newInterpreter()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	app := &cli.App{
		Name:  "neva",
		Usage: "Flow-based programming language",
		Commands: []*cli.Command{
			{
				Name:      "run",
				Usage:     "Run neva program from source code in interpreter mode",
				Args:      true,
				ArgsUsage: "Provide path to the executable package",
				Action: func(cCtx *cli.Context) error {
					args := cCtx.Args()
					path := strings.TrimSuffix(args.First(), "/main.neva")
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

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
