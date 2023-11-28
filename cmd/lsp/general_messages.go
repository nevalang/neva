package main

import (
	"context"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) Initialize(glspCtx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	result := protocol.InitializeResult{
		Capabilities: s.handler.CreateServerCapabilities(),
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    s.name,
			Version: &s.version,
		},
	}

	if params.RootPath == nil {
		glspCtx.Notify("neva/show_warning", "folder must be opened")
		return result, nil
	}

	// first indexing
	prog, problems, err := s.indexer.index(context.Background(), *params.RootPath)
	if err != nil {
		return nil, err
	}

	s.setState(prog, problems)

	return result, nil
}

// Initialized is called when vscode-extension is initialized.
// It spawns goroutines for sending indexing messages and warnings
func (s *Server) Initialized(glspCtx *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

// Shutdown closes channels so we don't miss any sent while vscode reloading.
// It shouldn't matter in production mode where vscode (re)launches server by itself.
// But is handy for development when server is spawned by user for debugging purposes.
func (s *Server) Shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (srv Server) Exit(glspContext *glsp.Context) error {
	return nil
}

func (srv Server) SetTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
