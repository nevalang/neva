package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/coder"
	"github.com/emil14/neva/internal/compiler/parser"
	"github.com/emil14/neva/internal/compiler/program"
	"github.com/emil14/neva/internal/compiler/translator"
	"github.com/emil14/neva/internal/compiler/validator"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/operators"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name: "neva",
		Commands: []*cli.Command{
			{
				Name: "build",
				Action: func(*cli.Context) error {
					ops := program.NewOperators()
					cmplr := compiler.MustNew(
						parser.MustNewYAML(),
						validator.MustNew(),
						translator.New(ops),
						coder.New(),
						ops,
					)

					dat, err := ioutil.ReadFile(`C:\projects\refactored-garbanzo\examples\arr.yml`)
					if err != nil {
						return err
					}

					bb, err := cmplr.Compile(dat)
					if err != nil {
						return err
					}

					return ioutil.WriteFile(
						`C:\projects\refactored-garbanzo\examples\arr.json`, bb, 0644,
					)
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
