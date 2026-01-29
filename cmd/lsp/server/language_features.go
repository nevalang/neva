package server

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) TextDocumentCompletion(
	glspCtx *glsp.Context,
	params *protocol.CompletionParams,
) (any, error) {
	s.logger.Info("TextDocumentCompletion")
	return []protocol.CompletionItem{}, nil
}
