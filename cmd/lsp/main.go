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

	h := &protocol.Handler{}
	s := Server{
		h:      h,
		logger: commonlog.GetLoggerf("%s.server", serverName),
	}

	// Base Protocol
	h.CancelRequest = nil
	h.Progress = nil

	// General Messages
	h.Initialize = s.initialize
	h.Initialized = s.initialized
	h.Shutdown = s.shutdown
	h.Exit = nil
	h.LogTrace = nil
	h.SetTrace = s.setTrace

	// Window
	h.WindowWorkDoneProgressCancel = nil

	// Workspace
	h.WorkspaceDidChangeWorkspaceFolders = nil
	h.WorkspaceDidChangeConfiguration = nil
	h.WorkspaceDidChangeWatchedFiles = nil
	h.WorkspaceSymbol = nil
	h.WorkspaceExecuteCommand = nil
	h.WorkspaceWillCreateFiles = nil
	h.WorkspaceDidCreateFiles = nil
	h.WorkspaceWillRenameFiles = nil
	h.WorkspaceDidRenameFiles = nil
	h.WorkspaceWillDeleteFiles = nil
	h.WorkspaceDidDeleteFiles = nil
	h.WorkspaceSemanticTokensRefresh = nil

	// Text Document Synchronization
	h.TextDocumentDidOpen = nil
	h.TextDocumentDidChange = nil
	h.TextDocumentWillSave = nil
	h.TextDocumentWillSaveWaitUntil = nil
	h.TextDocumentDidSave = nil
	h.TextDocumentDidClose = nil

	// Language features
	h.TextDocumentCompletion = nil
	h.CompletionItemResolve = nil
	h.TextDocumentHover = nil
	h.TextDocumentSignatureHelp = nil
	h.TextDocumentDeclaration = nil
	h.TextDocumentDefinition = nil
	h.TextDocumentTypeDefinition = nil
	h.TextDocumentImplementation = nil
	h.TextDocumentReferences = nil
	h.TextDocumentDocumentHighlight = nil
	h.TextDocumentDocumentSymbol = nil
	h.TextDocumentCodeAction = nil
	h.CodeActionResolve = nil
	h.TextDocumentCodeLens = nil
	h.CodeLensResolve = nil
	h.TextDocumentDocumentLink = nil
	h.DocumentLinkResolve = nil
	h.TextDocumentColor = nil
	h.TextDocumentColorPresentation = nil
	h.TextDocumentFormatting = nil
	h.TextDocumentRangeFormatting = nil
	h.TextDocumentOnTypeFormatting = nil
	h.TextDocumentRename = nil
	h.TextDocumentPrepareRename = nil
	h.TextDocumentFoldingRange = nil
	h.TextDocumentSelectionRange = nil
	h.TextDocumentPrepareCallHierarchy = nil
	h.CallHierarchyIncomingCalls = nil
	h.CallHierarchyOutgoingCalls = nil
	h.TextDocumentSemanticTokensFull = nil
	h.TextDocumentSemanticTokensFullDelta = nil
	h.TextDocumentSemanticTokensRange = nil
	h.TextDocumentLinkedEditingRange = nil
	h.TextDocumentMoniker = nil

	srv := server.NewServer(h, serverName, *isDebug)
	if err := srv.RunTCP(fmt.Sprintf("localhost:%d", *port)); err != nil {
		panic(err)
	}
}

type Server struct {
	h      *protocol.Handler // readonly
	logger commonlog.Logger
}

func (srv Server) initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	srv.logger.Info("initialize")

	return protocol.InitializeResult{
		Capabilities: handler.CreateServerCapabilities(),
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    serverName,
			Version: &version,
		},
	}, nil
}

func (srv Server) initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func (srv Server) shutdown(context *glsp.Context) error {
	fmt.Println("shutdown")
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (srv Server) setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
