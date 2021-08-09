package main

import (
	"log"
	"os"

	"github.com/emil14/stream/internal/runtime"
	"github.com/emil14/stream/internal/runtime/operators"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name: "stream",
		Commands: []*cli.Command{
			{
				Name: "run",
				Action: func(*cli.Context) error {
					// New().Run()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func New() runtime.Runtime {
	return runtime.Runtime{
		Operators: map[string]func(runtime.RuntimeIO) error{
			"*": operators.Mul,
		},
	}
}
