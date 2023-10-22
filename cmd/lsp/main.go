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

const serverName = "neva"

var (
	version = "0.0.1"
	handler protocol.Handler
)

func main() { //nolint:funlen
	isDebug := flag.Bool("debug", false, "-debug")
	flag.Parse()

	port := flag.Int("port", 9000, "-port")
	flag.Parse()

	commonlog.Configure(1, nil)

	handler = protocol.Handler{
		// Base Protocol
		CancelRequest: nil,
		Progress:      nil,

		// General Messages
		Initialize:  initialize,
		Initialized: initialized,
		Shutdown:    shutdown,
		Exit:        nil,
		LogTrace:    nil,
		SetTrace:    setTrace,

		// Window
		WindowWorkDoneProgressCancel: nil,

		// Workspace
		WorkspaceDidChangeWorkspaceFolders: nil,
		WorkspaceDidChangeConfiguration:    nil,
		WorkspaceDidChangeWatchedFiles:     nil,
		WorkspaceSymbol:                    nil,
		WorkspaceExecuteCommand:            nil,
		WorkspaceWillCreateFiles:           nil,
		WorkspaceDidCreateFiles:            nil,
		WorkspaceWillRenameFiles:           nil,
		WorkspaceDidRenameFiles:            nil,
		WorkspaceWillDeleteFiles:           nil,
		WorkspaceDidDeleteFiles:            nil,
		WorkspaceSemanticTokensRefresh:     nil,

		// Text Document Synchronization
		TextDocumentDidOpen:           nil,
		TextDocumentDidChange:         nil,
		TextDocumentWillSave:          nil,
		TextDocumentWillSaveWaitUntil: nil,
		TextDocumentDidSave:           nil,
		TextDocumentDidClose:          nil,

		// Language features
		TextDocumentCompletion:              nil,
		CompletionItemResolve:               nil,
		TextDocumentHover:                   nil,
		TextDocumentSignatureHelp:           nil,
		TextDocumentDeclaration:             nil,
		TextDocumentDefinition:              nil,
		TextDocumentTypeDefinition:          nil,
		TextDocumentImplementation:          nil,
		TextDocumentReferences:              nil,
		TextDocumentDocumentHighlight:       nil,
		TextDocumentDocumentSymbol:          nil,
		TextDocumentCodeAction:              nil,
		CodeActionResolve:                   nil,
		TextDocumentCodeLens:                nil,
		CodeLensResolve:                     nil,
		TextDocumentDocumentLink:            nil,
		DocumentLinkResolve:                 nil,
		TextDocumentColor:                   nil,
		TextDocumentColorPresentation:       nil,
		TextDocumentFormatting:              nil,
		TextDocumentRangeFormatting:         nil,
		TextDocumentOnTypeFormatting:        nil,
		TextDocumentRename:                  nil,
		TextDocumentPrepareRename:           nil,
		TextDocumentFoldingRange:            nil,
		TextDocumentSelectionRange:          nil,
		TextDocumentPrepareCallHierarchy:    nil,
		CallHierarchyIncomingCalls:          nil,
		CallHierarchyOutgoingCalls:          nil,
		TextDocumentSemanticTokensFull:      nil,
		TextDocumentSemanticTokensFullDelta: nil,
		TextDocumentSemanticTokensRange:     nil,
		TextDocumentLinkedEditingRange:      nil,
		TextDocumentMoniker:                 nil,
	}

	srv := server.NewServer(&handler, serverName, *isDebug)
	if err := srv.RunTCP(fmt.Sprintf("localhost:%d", *port)); err != nil {
		panic(err)
	}
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	logger := commonlog.GetLoggerf("%s.server", serverName)
	logger.Info("initialize")

	return protocol.InitializeResult{
		Capabilities: handler.CreateServerCapabilities(),
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
