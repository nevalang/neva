package server

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) Initialize(glspCtx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	s.workspacePath = *params.RootPath
	return protocol.InitializeResult{
		Capabilities: s.handler.CreateServerCapabilities(),
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    s.name,
			Version: &s.version,
		},
	}, nil
}

func (s *Server) Initialized(glspCtx *glsp.Context, params *protocol.InitializedParams) error {
	return s.indexAndNotifyProblems(glspCtx.Notify)
}

func (s *Server) Shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func (s Server) Exit(_ *glsp.Context) error {
	return nil
}

func (s Server) SetTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
