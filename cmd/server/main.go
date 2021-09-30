package main

import (
	"io/ioutil"
	"net/http"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/coder"
	"github.com/emil14/neva/internal/compiler/parser"
	"github.com/emil14/neva/internal/compiler/program"
	"github.com/emil14/neva/internal/compiler/translator"
	"github.com/emil14/neva/internal/compiler/validator"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/connector"
	"github.com/emil14/neva/internal/runtime/operators"
)

type Server struct {
	storage  storage.Storage
	compiler compiler.Compiler
	runtime  runtime.Runtime
}

func (s Server) handle() {
	bb, _ := ioutil.ReadFile("/home/emil14/projects/neva/examples/program/prog.yml")
}

func main() {
	http.HandleFunc("/program", h)
	http.ListenAndServe(":8090", nil)
}

func MustNew() Server {
	ops := program.NewOperators()
	return Server{
		compiler: compiler.MustNew(
			parser.MustNewYAML(),
			validator.New(),
			translator.New(ops),
			coder.New(),
			ops,
		),
		runtime: runtime.New(
			connector.MustNew(operators.New(), nil, nil),
		),
	}
}
