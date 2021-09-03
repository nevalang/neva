package main

import (
	"fmt"
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
	"github.com/emil14/neva/internal/runtime/decoder"
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
					ops := program.NewOperators()
					cmplr := compiler.MustNew(
						parser.MustNewYAML(),
						validator.New(),
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
			{
				Name: "run",
				Action: func(*cli.Context) error {
					bb, err := ioutil.ReadFile(`C:\projects\refactored-garbanzo\examples\arr.json`)
					if err != nil {
						return err
					}

					c := decoder.MustNewJSON()
					prog, err := c.Decode(bb)
					if err != nil {
						return err
					}

					ops := operators.New()
					r := runtime.New(ops)
					io, err := r.Run(prog)
					if err != nil {
						return err
					}

					in, err := io.In.Port("x")
					if err != nil {
						return err
					}

					outport, err := io.Out.Port("y")
					if err != nil {
						return err
					}

					in <- runtime.NewIntMsg(42)
					v := <-outport
					fmt.Println(v)

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
