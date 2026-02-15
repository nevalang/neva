package main

import (
	"flag"

	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp/server"

	lspServer "github.com/nevalang/neva/cmd/lsp/server"
	builder "github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/compiler/typesystem"
	"github.com/nevalang/neva/pkg/indexer"
)

func main() {
	const serverName = "neva"

	isDebug := flag.Bool("debug", false, "-debug")
	flag.Parse()

	loglvl := 1
	if *isDebug {
		loglvl = 2
	}

	commonlog.Configure(loglvl, nil)
	logger := commonlog.GetLoggerf("%s.server", serverName)

	p := parser.New()

	terminator := typesystem.Terminator{}
	checker := typesystem.MustNewSubtypeChecker(terminator)
	resolver := typesystem.MustNewResolver(typesystem.Validator{}, checker, terminator)
	builder := builder.MustNew(p)

	indexer := indexer.New(builder, p, analyzer.MustNew(resolver), logger)

	handler := lspServer.BuildHandler(logger, serverName, indexer)

	srv := server.NewServer(
		handler,
		serverName,
		*isDebug,
	)

	if *isDebug {
		if err := srv.RunTCP("localhost:6007"); err != nil {
			panic(err)
		}
	} else {
		if err := srv.RunStdio(); err != nil {
			panic(err)
		}
	}
}
