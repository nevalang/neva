package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emil14/respect/internal/compiler"
	"github.com/emil14/respect/internal/compiler/coder"
	"github.com/emil14/respect/internal/compiler/parser"
	cprog "github.com/emil14/respect/internal/compiler/program"
	"github.com/emil14/respect/internal/compiler/storage"
	"github.com/emil14/respect/internal/compiler/translator"
	"github.com/emil14/respect/internal/compiler/validator"
	"github.com/emil14/respect/internal/runtime"
	"github.com/emil14/respect/internal/runtime/connector"
	rprog "github.com/emil14/respect/internal/runtime/program"
	"github.com/emil14/respect/pkg/sdk"
)

func main() {
	srv := MustNew()
	sdkCtrl := sdk.NewDefaultApiController(srv)
	router := sdk.NewRouter(sdkCtrl)

	log.Println("listening http://localhost:8090")
	if err := http.ListenAndServe(":8090", router); err != nil {
		log.Fatalln(err)
	}
}

type Storage interface {
	PkgDescriptor(path string) (compiler.PkgDescriptor, error)
}

type Server struct {
	validator compiler.Validator
	compiler  compiler.Compiler
	runtime   runtime.Runtime
	caster    Caster
	storage   Storage
}

func (s Server) OnSend(msg runtime.Msg, from rprog.PortAddr) runtime.Msg {
	fmt.Println(msg, from)
	return msg
}

func (s Server) OnReceive(msg runtime.Msg, from rprog.PortAddr, to rprog.PortAddr) {
	fmt.Println(msg, from, to)
}

func MustNew() Server {
	ops := cprog.NewOperators()
	store := storage.MustNew("")
	v := validator.New()
	cmplr := compiler.MustNew(
		parser.MustNewYAML(),
		v,
		translator.New(ops),
		coder.New(),
		ops,
	)
	opspaths := map[string]runtime.Operator{
		"%": func(runtime.IO) error {
			return nil
		},
		"*": func(runtime.IO) error {
			return nil
		},
		"&&": func(runtime.IO) error {
			return nil
		},
		"||": func(runtime.IO) error {
			return nil
		},
		">": func(runtime.IO) error {
			return nil
		},
		"filter": func(runtime.IO) error {
			return nil
		},
	}

	s := Server{
		caster:    caster{},
		storage:   store,
		compiler:  cmplr,
		validator: v,
	}
	s.runtime = runtime.New(
		connector.MustNew(
			opspaths,
			s,
		),
	)

	return s
}
