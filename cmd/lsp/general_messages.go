package main

import (
	"context"
	"fmt"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s Server) Initialize(glpsCtx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	result := protocol.InitializeResult{
		Capabilities: s.handler.CreateServerCapabilities(),
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    s.name,
			Version: &s.version,
		},
	}

	if params.RootPath == nil {
		glpsCtx.Notify("neva/show_warning", "folder must be opened")
		return result, nil
	}

	prog, problems, err := s.indexer.index(context.Background(), *params.RootPath)
	if err != nil {
		return nil, fmt.Errorf("index: %w", err)
	}

	glpsCtx.Notify("neva/workdir_indexed", prog)
	if problems != "" {
		glpsCtx.Notify("neva/analyzer_message", prog)
	}

	return result, nil
}

func (srv Server) Initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func (srv Server) Shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (srv Server) SetTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
