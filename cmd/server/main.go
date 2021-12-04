package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emil14/neva/internal/compiler"
	checker "github.com/emil14/neva/internal/compiler/checker"
	"github.com/emil14/neva/internal/compiler/coder"
	"github.com/emil14/neva/internal/compiler/parser"
	"github.com/emil14/neva/internal/compiler/translator"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/connector"
	"github.com/emil14/neva/pkg/sdk"
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

type Server struct {
	compiler compiler.Compiler
	runtime  runtime.Runtime
	caster   Caster
}

func (s Server) OnSend(msg runtime.Msg, from runtime.PortAddr) runtime.Msg {
	fmt.Println(msg, from)
	return msg
}

func (s Server) OnReceive(msg runtime.Msg, from runtime.PortAddr, to runtime.PortAddr) {
	fmt.Println(msg, from, to)
}

func MustNew() Server {
	// opsIO := cprog.NewOperatorsIO()
	check := checker.New()

	cmplr := compiler.MustNew(
		parser.MustNewYAML(),
		check,
		translator.New(),
		coder.New(),
		nil,
		nil,
		// opsIO,
	)

	s := Server{
		caster:   caster{},
		compiler: cmplr,
	}

	s.runtime = runtime.New(
		connector.MustNew(nil),
	)

	return s
}
