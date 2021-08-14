package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/coder"
	"github.com/emil14/neva/internal/compiler/parser"
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
				Name: "compile",
				Action: func(*cli.Context) error {
					var (
						p    = parser.MustNewYAML()
						v    = validator.MustNew()
						t    = translator.New()
						c    = coder.MustNewJSON()
						comp = compiler.MustNew(p, v, t, c)
					)

					dat, err := ioutil.ReadFile(`C:\projects\refactored-garbanzo\examples\arr.yml`)
					if err != nil {
						return err
					}

					bb, err := comp.Compile(dat)
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
