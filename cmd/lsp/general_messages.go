package main

import (
	"context"

	"github.com/nevalang/neva/internal/compiler/sourcecode"
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

	// it's important to open these channels here and not when server is created
	// to avoid sending messages to closed connection, which can happen when vscode is reloading in development mode
	s.indexChan = make(chan sourcecode.Module)
	s.problemsChan = make(chan string)

	// first indexing
	go func() {
		prog, problems, err := s.indexer.index(context.Background(), *params.RootPath)
		if err != nil {
			s.logger.Errorf("indexer: %v", err.Error())
			return
		}
		s.indexChan <- prog
		s.problemsChan <- string(problems)
	}()

	return result, nil
}

// Initialized is called when vscode-extension is initialized.
// It spawns goroutines for sending indexing messages and warnings
func (s *Server) Initialized(glspCtx *glsp.Context, params *protocol.InitializedParams) error {
	go func() {
		for indexedMod := range s.indexChan {
			glspCtx.Notify("neva/workdir_indexed", indexedMod)
		}
	}()
	go func() {
		for problems := range s.problemsChan {
			glspCtx.Notify("neva/analyzer_message", problems)
		}
	}()
	return nil
}

// Shutdown closes channels so we don't miss any sent while vscode reloading.
// It shouldn't matter in production mode where vscode (re)launches server by itself.
// But is handy for development when server is spawned by user for debugging purposes.
func (s *Server) Shutdown(context *glsp.Context) error {
	close(s.indexChan)
	close(s.problemsChan)
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (srv Server) Exit(glspContext *glsp.Context) error {
	close(srv.indexChan)
	close(srv.problemsChan)
	return nil
}

func (srv Server) SetTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
