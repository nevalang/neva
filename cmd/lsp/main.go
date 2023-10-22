package main

import (
	"flag"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

func main() { //nolint:funlen
	const serverName = "neva"

	// cli args
	isDebug := flag.Bool("debug", false, "-debug")
	flag.Parse()
	port := flag.Int("port", 9000, "-port")
	flag.Parse()

	// logger
	verbosity := 1
	if *isDebug {
		verbosity = 2
	}
	commonlog.Configure(verbosity, nil)
	logger := commonlog.GetLoggerf("%s.server", serverName)

	// handler and server
	h := &protocol.Handler{}
	s := Server{
		handler: h,
		logger:  logger,
		parser:  parser.New(*isDebug),
	}

	// Base Protocol
	h.CancelRequest = func(context *glsp.Context, params *protocol.CancelParams) error {
		return nil
	}
	h.Progress = func(context *glsp.Context, params *protocol.ProgressParams) error {
		return nil
	}

	// General Messages
	h.Initialize = s.initialize
	h.Initialized = s.initialized
	h.Shutdown = s.shutdown
	h.Exit = func(context *glsp.Context) error {
		return nil
	}
	h.LogTrace = func(context *glsp.Context, params *protocol.LogTraceParams) error {
		return nil
	}
	h.SetTrace = s.setTrace

	// Window
	h.WindowWorkDoneProgressCancel = func(context *glsp.Context, params *protocol.WorkDoneProgressCancelParams) error {
		return nil
	}

	// Workspace
	h.WorkspaceDidChangeWorkspaceFolders = func(context *glsp.Context, params *protocol.DidChangeWorkspaceFoldersParams) error { //nolint:lll
		return nil
	}
	h.WorkspaceDidChangeConfiguration = func(context *glsp.Context, params *protocol.DidChangeConfigurationParams) error {
		return nil
	}
	h.WorkspaceDidChangeWatchedFiles = func(context *glsp.Context, params *protocol.DidChangeWatchedFilesParams) error {
		return nil
	}
	h.WorkspaceSymbol = func(context *glsp.Context, params *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) { //nolint:lll
		return nil, nil
	}
	h.WorkspaceExecuteCommand = func(context *glsp.Context, params *protocol.ExecuteCommandParams) (any, error) {
		return nil, nil
	}
	h.WorkspaceWillCreateFiles = func(context *glsp.Context, params *protocol.CreateFilesParams) (*protocol.WorkspaceEdit, error) { //nolint:lll
		return nil, nil
	}
	h.WorkspaceDidCreateFiles = func(context *glsp.Context, params *protocol.CreateFilesParams) error {
		return nil
	}
	h.WorkspaceWillRenameFiles = func(context *glsp.Context, params *protocol.RenameFilesParams) (*protocol.WorkspaceEdit, error) { //nolint:lll
		return nil, nil
	}
	h.WorkspaceDidRenameFiles = func(context *glsp.Context, params *protocol.RenameFilesParams) error {
		return nil
	}
	h.WorkspaceWillDeleteFiles = func(context *glsp.Context, params *protocol.DeleteFilesParams) (*protocol.WorkspaceEdit, error) { //nolint:lll
		return nil, nil
	}
	h.WorkspaceDidDeleteFiles = func(context *glsp.Context, params *protocol.DeleteFilesParams) error {
		return nil
	}
	h.WorkspaceSemanticTokensRefresh = func(context *glsp.Context) error {
		return nil
	}

	// Text Document Synchronization
	h.TextDocumentDidOpen = s.TextDocumentDidOpen
	h.TextDocumentDidChange = func(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
		return nil
	}
	h.TextDocumentWillSave = func(context *glsp.Context, params *protocol.WillSaveTextDocumentParams) error {
		return nil
	}
	h.TextDocumentWillSaveWaitUntil = func(context *glsp.Context, params *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) { //nolint:lll
		return nil, nil
	}
	h.TextDocumentDidSave = func(context *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
		return nil
	}
	h.TextDocumentDidClose = func(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
		return nil
	}

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
	name, version string
	parser        parser.Parser
	handler       *protocol.Handler // readonly
	logger        commonlog.Logger
}

func (srv Server) initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	return protocol.InitializeResult{
		Capabilities: srv.handler.CreateServerCapabilities(),
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    srv.name,
			Version: &srv.version,
		},
	}, nil
}

func (srv Server) initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func (srv Server) shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (srv Server) setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
