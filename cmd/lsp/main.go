package main

import (
	"flag"
	"sync"

	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp/server"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/parser"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/pkg/typesystem"
)

type Server struct {
	name, version string

	handler *Handler
	logger  commonlog.Logger
	indexer Indexer

	mu    *sync.Mutex
	state *State
}

type State struct {
	mod     src.Module
	problem string
}

func main() {
	const serverName = "neva"

	isDebug := flag.Bool("debug", false, "-debug")
	flag.Parse()

	verbosity := 1
	if *isDebug {
		verbosity = 2
	}
	commonlog.Configure(verbosity, nil)
	logger := commonlog.GetLoggerf("%s.server", serverName)

	p := parser.MustNew(*isDebug)

	terminator := typesystem.Terminator{}
	checker := typesystem.MustNewSubtypeChecker(terminator)
	resolver := typesystem.MustNewResolver(typesystem.Validator{}, checker, terminator)
	builder := builder.MustNew(
		"/Users/emil/projects/neva/std",
		"/Users/emil/projects/neva/third_party/",
		p,
	)

	indexer := Indexer{
		builder:  builder,
		parser:   p,
		analyzer: analyzer.MustNew(resolver),
	}

	handler := buildHandler(logger, serverName, indexer)

	srv := server.NewServer(
		handler,
		serverName,
		*isDebug,
	)

	if err := srv.RunStdio(); err != nil {
		panic(err)
	}
}
