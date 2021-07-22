package main

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name: "fbp",
		Commands: []*cli.Command{
			{
				Name:   "run",
				Action: run,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
