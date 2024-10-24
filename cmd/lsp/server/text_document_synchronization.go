package server

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) TextDocumentDidOpen(
	glspCtx *glsp.Context,
	params *protocol.DidOpenTextDocumentParams,
) error {
	s.activeFileMutex.Lock()
	s.activeFile = params.TextDocument.URI
	s.activeFileMutex.Unlock()
	return nil
}

func (s *Server) TextDocumentDidChange(
	glspCtx *glsp.Context,
	params *protocol.DidChangeTextDocumentParams,
) error {
	s.activeFileMutex.Lock()
	s.activeFile = params.TextDocument.URI
	s.activeFileMutex.Unlock()
	return nil
}

func (s *Server) TextDocumentDidSave(
	glspCtx *glsp.Context,
	params *protocol.DidSaveTextDocumentParams,
) error {
	s.logger.Info("TextDocumentDidSave")
	return s.indexAndNotifyProblems(glspCtx.Notify)
}
