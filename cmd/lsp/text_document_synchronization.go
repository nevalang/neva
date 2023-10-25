package main

import (
	"context"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Message struct {
	Text string `json:"text"`
}

func (s Server) TextDocumentDidOpen(glpsCtx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	parsedFile, err := s.parser.ParseFile(context.Background(), []byte(params.TextDocument.Text))
	if err != nil {
		return err
	}

	glpsCtx.Notify("neva/file_parsed", parsedFile)

	return nil
}
