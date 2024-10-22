package server

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) TextDocumentDidChange(
	glspCtx *glsp.Context,
	params *protocol.DidChangeTextDocumentParams,
) error {
	return nil
}

func (s *Server) TextDocumentDidSave(
	glspCtx *glsp.Context,
	params *protocol.DidSaveTextDocumentParams,
) error {
	s.logger.Info("TextDocumentDidSave")
	return s.indexAndNotifyProblems(glspCtx.Notify)
}
