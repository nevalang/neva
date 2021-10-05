package main

import (
	"fmt"
	"net/http"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/coder"
	"github.com/emil14/neva/internal/compiler/parser"
	"github.com/emil14/neva/internal/compiler/program"
	"github.com/emil14/neva/internal/compiler/storage"
	"github.com/emil14/neva/internal/compiler/translator"
	"github.com/emil14/neva/internal/compiler/validator"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/connector"
	"github.com/emil14/neva/internal/runtime/operators"
)

type Server struct {
	compiler compiler.Compiler
	runtime  runtime.Runtime
}

func (s Server) handle(http.ResponseWriter, *http.Request) {
	prog, err := s.compiler.PreCompile("/home/emil14/projects/neva/examples/program/prog.yml")
	if err != nil {
		panic(err)
	}

	_, err = s.runtime.Run(prog)
	if err != nil {
		panic(err)
	}
}

func (s Server) onSend(msg runtime.Msg, from runtime.PortAddr, to []runtime.PortAddr) runtime.Msg {
	fmt.Println(msg, from, to)
	return msg
}

func (s Server) onReceive(msg runtime.Msg, from runtime.PortAddr, to runtime.PortAddr) {
	fmt.Println(msg, from, to)
}

func main() {
	s := MustNew()
	http.HandleFunc("/", s.handle)
	http.ListenAndServe(":8090", nil)
}

func MustNew() Server {
	s := Server{}
	ops := program.NewOperators()
	s.compiler = compiler.MustNew(
		parser.MustNewYAML(),
		validator.New(),
		translator.New(ops),
		coder.New(),
		storage.MustNew(""),
		ops,
	)
	s.runtime = runtime.New(connector.MustNew(operators.New(), s.onSend, s.onReceive))
	return s
}
