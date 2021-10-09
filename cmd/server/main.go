package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/coder"
	"github.com/emil14/neva/internal/compiler/parser"
	cprog "github.com/emil14/neva/internal/compiler/program"
	"github.com/emil14/neva/internal/compiler/storage"
	"github.com/emil14/neva/internal/compiler/translator"
	"github.com/emil14/neva/internal/compiler/validator"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/connector"
	"github.com/emil14/neva/internal/runtime/operators"
	rprog "github.com/emil14/neva/internal/runtime/program"
)

type Server struct {
	compiler compiler.Compiler
	runtime  runtime.Runtime
}

func (s Server) handle(w http.ResponseWriter, r *http.Request) {
	p, err := filepath.Abs("../../examples/program/pkg.yml")
	if err != nil {
		log.Println(err)
		return
	}

	prog, cprog, err := s.compiler.PreCompile(p)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(cprog)

	io, err := s.runtime.Run(prog)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(io)

}

func (s Server) onSend(msg runtime.Msg, from rprog.PortAddr, to []rprog.PortAddr) runtime.Msg {
	fmt.Println(msg, from, to)
	return msg
}

func (s Server) onReceive(msg runtime.Msg, from rprog.PortAddr, to rprog.PortAddr) {
	fmt.Println(msg, from, to)
}

func main() {
	s := MustNew()
	http.HandleFunc("/", s.handle)
	http.ListenAndServe(":8090", nil)
}

func MustNew() Server {
	s := Server{}
	ops := cprog.NewOperators()
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
