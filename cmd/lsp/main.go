package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/nevalang/neva/pkg/typesystem"
)

type Server struct {
	name, version string

	handler *protocol.Handler // readonly
	logger  commonlog.Logger
	indexer Indexer

	state src.Program // TODO protect with mutex (or replace with channel)
}

type Indexer struct {
	builder  builder.Builder
	frontend compiler.FrontEnd
}

func (p Indexer) process(ctx context.Context, path string) (src.Program, error) {
	raw, err := p.builder.Build(ctx, path)
	if err != nil {
		return src.Program{}, fmt.Errorf("builder: %w", err)
	}

	prog, err := p.frontend.Process(ctx, raw, "")
	if err != nil {
		return src.Program{}, fmt.Errorf("frontend: %w", err)
	}

	return prog, nil
}

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

	// compiler
	terminator := typesystem.Terminator{}
	checker := typesystem.MustNewSubtypeChecker(terminator)
	resolver := typesystem.MustNewResolver(typesystem.Validator{}, checker, terminator)
	builder := builder.MustNew("/Users/emil/projects/neva/std")
	frontend := compiler.NewFrontEnd(parser.MustNew(*isDebug), analyzer.MustNew(resolver))

	// handler and server
	h := &protocol.Handler{}
	s := Server{
		handler: h,
		logger:  logger,
		name:    serverName,
		version: "0.0.1",
		indexer: Indexer{
			builder:  builder,
			frontend: frontend,
		},
	}

	// Base Protocol
	h.CancelRequest = func(context *glsp.Context, params *protocol.CancelParams) error {
		return nil
	}
	h.Progress = func(context *glsp.Context, params *protocol.ProgressParams) error {
		return nil
	}

	// General Messages
	h.Initialize = s.Initialize
	h.Initialized = s.Initialized
	h.Shutdown = s.Shutdown
	h.Exit = func(context *glsp.Context) error {
		return nil
	}
	h.LogTrace = func(context *glsp.Context, params *protocol.LogTraceParams) error {
		return nil
	}
	h.SetTrace = s.SetTrace

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
