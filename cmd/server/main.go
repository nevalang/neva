package main

import (
	"context"
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

	"github.com/emil14/neva/pkg/sdk"
)

type Server struct {
	compiler compiler.Compiler
	runtime  runtime.Runtime
}

func (s Server) ProgramGet(context.Context) (sdk.ImplResponse, error) {
	p, err := filepath.Abs("../../examples/program/pkg.yml")
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	prog, cprog, err := s.compiler.PreCompile(p)
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	log.Println(cprog)

	io, err := s.runtime.Run(prog)
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	log.Println(io)

	return sdk.ImplResponse{
		Code: 0,
		Body: nil,
	}, nil
}

func (s Server) onSend(msg runtime.Msg, from rprog.PortAddr, to []rprog.PortAddr) runtime.Msg {
	fmt.Println(msg, from, to)
	return msg
}

func (s Server) onReceive(msg runtime.Msg, from rprog.PortAddr, to rprog.PortAddr) {
	fmt.Println(msg, from, to)
}

func main() {
	srv := MustNew()
	ctrl := sdk.NewDefaultApiController(srv)
	r := sdk.NewRouter(ctrl)
	http.ListenAndServe(":8090", r)
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
