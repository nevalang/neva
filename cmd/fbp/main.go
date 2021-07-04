package main

import (
	"errors"
	"fbp/internal/api"
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "fbp",
		Commands: []*cli.Command{
			{
				Name: "start",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name: "port",
					},
				},
				Action: func(ctx *cli.Context) error {
					port, ok := ctx.Value("port").(int)
					if !ok {
						return errors.New("port not int")
					}
					return api.ListenAndServe(fmt.Sprintf(":%d", port))
				},
			},
			// {
			// 	Name:   "load",
			// 	Action: load,
			// },
			{
				Name:   "parse",
				Action: parse,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
