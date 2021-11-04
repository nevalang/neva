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
	"github.com/emil14/respect/internal/core"
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
	compiler compiler.Compiler
	runtime  runtime.Runtime
	caster   Caster
	storage  Storage
}

func (s Server) OnSend(msg core.Msg, from rprog.PortAddr) core.Msg {
	fmt.Println(msg, from)
	return msg
}

func (s Server) OnReceive(msg core.Msg, from rprog.PortAddr, to rprog.PortAddr) {
	fmt.Println(msg, from, to)
}

func MustNew() Server {
	compilerOps := cprog.NewOperators()
	compilerValidator := validator.New()
	cmplr := compiler.MustNew(
		parser.MustNewYAML(),
		compilerValidator,
		translator.New(compilerOps),
		coder.New(),
		compilerOps,
	)

	// fakeRuntimeOps := map[string]runtime.OperatorFunc{
	// 	"%": func(core.IO) error {
	// 		return nil
	// 	},
	// 	"*": func(core.IO) error {
	// 		return nil
	// 	},
	// 	"&&": func(core.IO) error {
	// 		return nil
	// 	},
	// 	"||": func(core.IO) error {
	// 		return nil
	// 	},
	// 	">": func(core.IO) error {
	// 		return nil
	// 	},
	// 	"filter": func(core.IO) error {
	// 		return nil
	// 	},
	// }

	s := Server{
		caster:   caster{},
		storage:  storage.MustNew(""),
		compiler: cmplr,
	}

	s.runtime = runtime.New(connector.MustNew(s))

	return s
}
