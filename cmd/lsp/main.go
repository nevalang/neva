package main

import (
	"flag"
	"fmt"

	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const serverName = "neva:lsp"

var (
	version = "0.0.1"
	handler protocol.Handler
)

func main() {
	isDebug := flag.Bool("debug", false, "-debug")
	flag.Parse()

	port := flag.Int("port", 9000, "-port")
	flag.Parse()

	commonlog.Configure(1, nil)

	handler = protocol.Handler{
		Initialize:  initialize,
		Initialized: initialized,
		Shutdown:    shutdown,
		SetTrace:    setTrace,
	}

	srv := server.NewServer(&handler, serverName, *isDebug)
	if err := srv.RunTCP(fmt.Sprintf("localhost:%d", *port)); err != nil {
		panic(err)
	}
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	fmt.Println("initialize")

	capabilities := handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    serverName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	fmt.Println("shutdown")
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
