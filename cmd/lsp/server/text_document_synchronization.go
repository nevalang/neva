package server

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) TextDocumentDidChange(
	glspCtx *glsp.Context,
	params *protocol.DidChangeTextDocumentParams,
) error {
	s.logger.Info("TextDocumentDidChange")
	return s.indexAndNotifyProblems(glspCtx.Notify)
}
