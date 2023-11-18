package main

import (
	"context"
	"fmt"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s Server) Initialize(glspCtx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	fmt.Println("===Initialize===")

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

	// first indexation
	go func() {
		prog, problems, err := s.indexer.index(context.Background(), *params.RootPath)
		if err != nil {
			s.logger.Errorf("indexer: %v", err.Error())
			return
		}

		s.indexChan <- prog
		s.problems <- string(problems)
	}()

	return result, nil
}

// Initialized is called when vscode-extension is initialized.
// It spawns goroutines for sending indexing messages and warnings
// Note that this methods only works correctly if any time vscode reloaded it relaunches language-server.
func (srv Server) Initialized(glspCtx *glsp.Context, params *protocol.InitializedParams) error {
	go func() {
		for indexedMod := range srv.indexChan {
			glspCtx.Notify("neva/workdir_indexed", indexedMod)
		}
	}()
	go func() {
		for problems := range srv.problems {
			glspCtx.Notify("neva/analyzer_message", problems)
		}
	}()
	return nil
}

func (srv Server) Shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (srv Server) Exit(glspContext *glsp.Context) error {
	close(srv.indexChan)
	close(srv.problems)
	return nil
}

func (srv Server) SetTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
