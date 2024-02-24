package main

import (
	"flag"

	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp/server"

	"github.com/nevalang/neva/cmd/lsp/indexer"
	lspServer "github.com/nevalang/neva/cmd/lsp/server"
	builder "github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/desugarer"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/pkg"
	"github.com/nevalang/neva/pkg/typesystem"
)

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

	p := parser.New(*isDebug)

	terminator := typesystem.Terminator{}
	checker := typesystem.MustNewSubtypeChecker(terminator)
	resolver := typesystem.MustNewResolver(typesystem.Validator{}, checker, terminator)
	builder := builder.New(
		"/Users/emil/projects/neva/std",
		"/Users/emil/projects/neva/third_party/",
		p,
	)

	indexer := indexer.New(
		builder,
		p,
		desugarer.Desugarer{},
		analyzer.MustNew(pkg.Version, resolver),
	)

	handler := lspServer.BuildHandler(logger, serverName, indexer)

	srv := server.NewServer(
		handler,
		serverName,
		*isDebug,
	)

	if err := srv.RunStdio(); err != nil {
		panic(err)
	}
}
