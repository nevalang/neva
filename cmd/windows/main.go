package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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
	rprog "github.com/emil14/neva/internal/runtime/program"
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

func (s Server) ProgramGet(ctx context.Context, path string) (sdk.ImplResponse, error) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	if path == "" {
		path = "examples/program/pkg.yml"
	}
	p := filepath.Join(pwd, "../../", path)

	pkgd, err := s.storage.PkgDescriptor(p)
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	rprog, cprog, err := s.compiler.BuildProgram(pkgd)
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	if _, err = s.runtime.Run(rprog); err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	casted, err := s.caster.toSDK(cprog)
	if err != nil {
		return sdk.ImplResponse{}, err
	}

	return sdk.ImplResponse{
		Code: 200,
		Body: casted,
	}, nil
}

func (s Server) ProgramPatch(context.Context, string, sdk.Program) (sdk.ImplResponse, error) {
	return sdk.ImplResponse{
		Code: 0,
		Body: nil,
	}, nil
}

func (s Server) ProgramPost(ctx context.Context, path string, prog sdk.Program) (sdk.ImplResponse, error) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	pkgd, err := s.storage.PkgDescriptor(filepath.Join(pwd, "../../", path))
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	rprog, cprog, err := s.compiler.BuildProgram(pkgd)
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	if _, err = s.runtime.Run(rprog); err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	casted, err := s.caster.toSDK(cprog)
	if err != nil {
		return sdk.ImplResponse{}, err
	}

	return sdk.ImplResponse{
		Code: 200,
		Body: casted,
	}, nil
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
