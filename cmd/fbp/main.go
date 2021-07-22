package main

import (
	"log"
	"os"

	"github.com/emil14/refactored-garbanzo/internal/parser"
	"github.com/emil14/refactored-garbanzo/internal/translator"
	cli "github.com/urfave/cli/v2"
)

var (
	t = translator.New()
	v = parser.NewValidator()
	p = parser.NewYAMLParser()
)

func main() {
	app := cli.App{
		Name: "fbp",
		Commands: []*cli.Command{
			{
				Name:   "run",
				Action: run,
			},
			{
				Name:   "check",
				Action: check,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
