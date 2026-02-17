package main

import (
	"flag"

	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp/server"

	lspServer "github.com/nevalang/neva/cmd/lsp/server"
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

	indexer := indexer.MustNewDefault(logger)

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
