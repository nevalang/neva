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
	"github.com/emil14/neva/internal/runtime/loader"
	rprog "github.com/emil14/neva/internal/runtime/program"
	"github.com/emil14/neva/pkg/sdk"
)

type Server struct {
	compiler compiler.Compiler
	runtime  runtime.Runtime
	caster   Caster
}

func (s Server) ProgramGet(ctx context.Context, path string) (sdk.ImplResponse, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	p, err := filepath.Abs(filepath.Join(dir, "examples/program/pkg.yml"))
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	prog, cprog, err := s.compiler.PreCompile(p)
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	io, err := s.runtime.Run(prog)
	if err != nil {
		log.Println(err)
		return sdk.ImplResponse{}, err
	}

	log.Println(io)

	casted, err := s.caster.CastProgram(cprog)
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

func (s Server) ProgramPost(context.Context, string, sdk.Program) (sdk.ImplResponse, error) {
	return sdk.ImplResponse{
		Code: 0,
		Body: nil,
	}, nil
}

func (s Server) OnSend(msg runtime.Msg, from rprog.PortAddr, to []rprog.PortAddr) runtime.Msg {
	fmt.Println(msg, from, to)
	return msg
}

func (s Server) OnReceive(msg runtime.Msg, from rprog.PortAddr, to rprog.PortAddr) {
	fmt.Println(msg, from, to)
}

func main() {
	srv := MustNew()
	ctrl := sdk.NewDefaultApiController(srv)
	r := sdk.NewRouter(ctrl)

	log.Println("listening http://localhost:8090")
	if err := http.ListenAndServe(":8090", r); err != nil {
		log.Fatalln(err)
	}
}

func MustNew() Server {
	s := Server{
		caster: caster{},
	}

	compilerOps := cprog.NewOperators()

	s.compiler = compiler.MustNew(
		parser.MustNewYAML(),
		validator.New(),
		translator.New(compilerOps),
		coder.New(),
		storage.MustNew(""),
		compilerOps,
	)

	// TODO: move this step inside runtime's Run method
	// otherwise any fbp-probram will have all the operators loaded
	// which makes all the plugin-related effort useless
	opspaths, err := loader.Load(map[string]loader.Params{
		"*": {
			PluginPath:     "plugins/mul.so",
			ExportedEntity: "Mul",
		},
		"&&": {
			PluginPath:     "plugins/and.so",
			ExportedEntity: "And",
		},
		"||": {
			PluginPath:     "plugins/or.so",
			ExportedEntity: "Or",
		},
		">": {
			PluginPath:     "plugins/more.so",
			ExportedEntity: "More",
		},
		"filter": {
			PluginPath:     "plugins/filter.so",
			ExportedEntity: "Filter",
		},
	})
	if err != nil {
		panic(err)
	}

	s.runtime = runtime.New(
		connector.MustNew(
			opspaths,
			s,
		),
	)

	return s
}

type Caster interface {
	CastProgram(cprog.Program) (sdk.Program, error)
}

type caster struct{}

func (c caster) CastProgram(from cprog.Program) (sdk.Program, error) {
	cc, err := c.castComponents(from.Components)
	if err != nil {
		return sdk.Program{}, err
	}
	return sdk.Program{
		Scope: cc,
		Root:  from.Root,
	}, nil
}

func (c caster) castComponents(from map[string]cprog.Component) (map[string]sdk.Component, error) {
	r := map[string]sdk.Component{}
	for k, v := range from {
		cmpnt, err := c.castComponent(v)
		if err != nil {
			return nil, err
		}
		r[k] = cmpnt
	}
	return r, nil
}

func (c caster) castComponent(from cprog.Component) (sdk.Component, error) {
	if _, ok := from.(cprog.Operator); ok {
		return sdk.Component{
			Io: c.castIO(from.Interface()),
		}, nil
	}

	mod, ok := from.(cprog.Module)
	if !ok {
		return sdk.Component{}, fmt.Errorf("casterr: unknown component type")
	}

	return sdk.Component{
		Io:      c.castIO(mod.Interface()),
		Workers: mod.Workers,
		Deps:    c.castDeps(mod.Deps),
		Het:     c.castNet(mod.Net),
	}, nil
}

func (c caster) castDeps(from map[string]cprog.IO) map[string]sdk.Io {
	r := map[string]sdk.Io{}
	for k, v := range from {
		r[k] = c.castIO(v)
	}
	return r
}

func (c caster) castNet(net cprog.OutgoingConnections) []sdk.Connection {
	r := make([]sdk.Connection, 0, len(net))
	for from, to := range net {
		for rcvr := range to {
			r = append(r, c.sdkConnection(from, rcvr))
		}
	}
	return r
}

func (c caster) sdkConnection(from, to cprog.PortAddr) sdk.Connection {
	return sdk.Connection{
		From: c.sdkPortAddr(from),
		To:   c.sdkPortAddr(to),
	}
}

func (c caster) sdkPortAddr(from cprog.PortAddr) sdk.PortAddr {
	return sdk.PortAddr{
		Node: from.Node,
		Idx:  int32(from.Idx),
		Port: from.Port,
	}
}

func (c caster) castIO(from cprog.IO) sdk.Io {
	return sdk.Io{
		In:  c.castPorts(from.In),
		Out: c.castPorts(from.Out),
	}
}

func (c caster) castPorts(from cprog.Ports) map[string]string {
	to := make(map[string]string, len(from))
	for name, typ := range from {
		to[name] = typ.String()
	}
	return to
}
