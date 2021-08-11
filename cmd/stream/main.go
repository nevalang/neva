package main

import (
	"fmt"
	"log"
	"os"

	"github.com/emil14/stream/internal/compiler"
	"github.com/emil14/stream/internal/compiler/coder"
	"github.com/emil14/stream/internal/compiler/parser"
	"github.com/emil14/stream/internal/compiler/translator"
	"github.com/emil14/stream/internal/compiler/validator"
	"github.com/emil14/stream/internal/runtime"
	"github.com/emil14/stream/internal/runtime/operators"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name: "stream",
		Commands: []*cli.Command{
			{
				Name: "compile",
				Action: func(*cli.Context) error {
					var (
						p    = parser.MustNewYAML()
						v    = validator.MustNew()
						t    = translator.MustNew()
						c    = coder.MustNewJSON()
						comp = compiler.MustNew(p, v, t, c)
					)

					_, err := comp.Compile([]byte(`
deps:
	"*":
	in:
		nums[]: int
	out:
		mul: int

in:
	x: int
out:
	y: int

workers:
	multi: "*"

net:
	in:
	x:
		multi:
		- nums[0]
		- nums[1]
	multi:
	mul:
		out: [y]
`))

					fmt.Println(string([]byte(`
deps:
	"*":
	in:
		nums[]: int
	out:
		mul: int

in:
	x: int
out:
	y: int

workers:
	multi: "*"

net:
	in:
	x:
		multi:
		- nums[0]
		- nums[1]
	multi:
	mul:
		out: [y]
`)))

					return err
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
		Operators: operators.New(),
	}
}
