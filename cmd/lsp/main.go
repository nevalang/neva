package main

import (
	"flag"

	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/parser"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/pkg/typesystem"
)

type Server struct {
	name, version string

	handler *protocol.Handler // readonly
	logger  commonlog.Logger
	indexer Indexer

	indexChan    chan src.Module
	problemsChan chan string
}

func main() {
	const serverName = "neva"

	isDebug := flag.Bool("debug", false, "-debug")
	flag.Parse()

	// verbosity := 1
	// if *isDebug {
	// 	verbosity = 2
	// }
	commonlog.Configure(0, nil)
	logger := commonlog.GetLoggerf("%s.server", serverName)

	p := parser.MustNew(*isDebug)

	terminator := typesystem.Terminator{}
	checker := typesystem.MustNewSubtypeChecker(terminator)
	resolver := typesystem.MustNewResolver(typesystem.Validator{}, checker, terminator)
	builder := builder.MustNew("/Users/emil/projects/neva/std", "/Users/emil/projects/neva/third_party/", p)

	indexer := Indexer{
		builder:  builder,
		parser:   p,
		analyzer: analyzer.MustNew(resolver),
	}

	if err := server.NewServer(
		buildHandler(logger, serverName, indexer),
		serverName,
		*isDebug,
	).RunStdio(); err != nil {
		panic(err)
	}
}
