package main

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Handler struct {
	*protocol.Handler

	FooBar func(glspCtx *glsp.Context, params FooBarParams) (any, error)
}

func (h Handler) Handle(glspCtx *glsp.Context) (response any, validMethod bool, validParams bool, err error) {
	if !h.IsInitialized() && (glspCtx.Method != protocol.MethodInitialize) {
		return nil, true, true, errors.New("server not initialized")
	}

	if glspCtx.Method == "foobar" {
		var params FooBarParams
		if err := json.Unmarshal(glspCtx.Params, &params); err != nil {
			return nil, true, false, err
		}

		resp, err := h.FooBar(glspCtx, params)
		if err != nil {
			return nil, true, true, err
		}

		return resp, true, true, nil
	}

	return h.Handle(glspCtx)
}

//nolint:lll,funlen
func buildHandler(logger commonlog.Logger, serverName string, indexer Indexer) *Handler {
	h := &Handler{
		Handler: &protocol.Handler{},
	}

	s := Server{
		handler:             h,
		logger:              logger,
		name:                serverName,
		version:             "0.0.1",
		indexer:             indexer,
		mu:                  &sync.Mutex{},
		indexedProgramState: nil,
	}

	// Basic
	h.CancelRequest = func(context *glsp.Context, params *protocol.CancelParams) error {
		return nil
	}
	h.Progress = func(context *glsp.Context, params *protocol.ProgressParams) error {
		return nil
	}

	// Lifetime
	h.Initialize = s.Initialize
	h.Initialized = s.Initialized
	h.Shutdown = s.Shutdown
	h.Exit = s.Exit
	h.LogTrace = func(context *glsp.Context, params *protocol.LogTraceParams) error {
		return nil
	}
	h.SetTrace = s.SetTrace

	// Custom handlers
	h.FooBar = s.FooBar

	// Rest...
	h.WindowWorkDoneProgressCancel = func(context *glsp.Context, params *protocol.WorkDoneProgressCancelParams) error {
		return nil
	}

	h.WorkspaceDidChangeWorkspaceFolders = func(context *glsp.Context, params *protocol.DidChangeWorkspaceFoldersParams) error {
		return nil
	}
	h.WorkspaceDidChangeConfiguration = func(context *glsp.Context, params *protocol.DidChangeConfigurationParams) error {
		return nil
	}
	h.WorkspaceDidChangeWatchedFiles = func(context *glsp.Context, params *protocol.DidChangeWatchedFilesParams) error {
		return nil
	}
	h.WorkspaceSymbol = func(context *glsp.Context, params *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
		return nil, nil
	}
	h.WorkspaceExecuteCommand = func(context *glsp.Context, params *protocol.ExecuteCommandParams) (any, error) {
		return nil, nil
	}
	h.WorkspaceWillCreateFiles = func(context *glsp.Context, params *protocol.CreateFilesParams) (*protocol.WorkspaceEdit, error) {
		return nil, nil
	}
	h.WorkspaceDidCreateFiles = func(context *glsp.Context, params *protocol.CreateFilesParams) error {
		return nil
	}
	h.WorkspaceWillRenameFiles = func(context *glsp.Context, params *protocol.RenameFilesParams) (*protocol.WorkspaceEdit, error) {
		return nil, nil
	}
	h.WorkspaceDidRenameFiles = func(context *glsp.Context, params *protocol.RenameFilesParams) error {
		return nil
	}
	h.WorkspaceWillDeleteFiles = func(context *glsp.Context, params *protocol.DeleteFilesParams) (*protocol.WorkspaceEdit, error) {
		return nil, nil
	}
	h.WorkspaceDidDeleteFiles = func(context *glsp.Context, params *protocol.DeleteFilesParams) error {
		return nil
	}
	h.WorkspaceSemanticTokensRefresh = func(context *glsp.Context) error {
		return nil
	}

	h.TextDocumentDidOpen = s.TextDocumentDidOpen
	h.TextDocumentDidChange = func(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
		return nil
	}
	h.TextDocumentWillSave = func(context *glsp.Context, params *protocol.WillSaveTextDocumentParams) error {
		return nil
	}
	h.TextDocumentWillSaveWaitUntil = func(context *glsp.Context, params *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
		return nil, nil
	}
	h.TextDocumentDidSave = func(context *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
		return nil
	}
	h.TextDocumentDidClose = func(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
		return nil
	}

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

	return h
}
