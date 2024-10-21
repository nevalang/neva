package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/nevalang/neva/cmd/lsp/indexer"
	"github.com/nevalang/neva/internal/compiler"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

type Server struct {
	workspacePath string
	name, version string

	handler *Handler
	logger  commonlog.Logger
	indexer indexer.Indexer

	mu    *sync.Mutex
	index *src.Build
}

// setState allows to update state in a thread-safe manner.
func (s *Server) updateIndex(build src.Build) {
	s.mu.Lock()
	s.index = &build
	s.mu.Unlock()
}

func (s *Server) indexAndNotifyProblems(notify glsp.NotifyFunc) error {
	build, err := s.indexer.FullIndex(context.Background(), s.workspacePath)

	s.updateIndex(build)

	if err == nil {
		// clear problems
		notify(
			protocol.ServerTextDocumentPublishDiagnostics,
			protocol.PublishDiagnosticsParams{
				Diagnostics: []protocol.Diagnostic{},
			},
		)
		s.logger.Info("full index without problems, sent empty diagnostics")
		return nil
	}

	notify(
		protocol.ServerTextDocumentPublishDiagnostics,
		s.createDiagnostics(*err),
	)

	s.logger.Info("diagnostic sent: " + err.Error())

	return nil
}

func (s *Server) createDiagnostics(compilerErr compiler.Error) protocol.PublishDiagnosticsParams {
	var uri string
	if compilerErr.Location != nil {
		uri = fmt.Sprintf(
			"%s/%s/%s",
			s.workspacePath,
			compilerErr.Location.PkgName,
			compilerErr.Location.FileName+".neva",
		)
	}

	var startStopRange protocol.Range
	if compilerErr.Range != nil {
		startStopRange = protocol.Range{
			Start: protocol.Position{
				Line:      uint32(compilerErr.Range.Start.Line),
				Character: uint32(compilerErr.Range.Start.Column),
			},
			End: protocol.Position{
				Line:      uint32(compilerErr.Range.Stop.Line),
				Character: uint32(compilerErr.Range.Stop.Column),
			},
		}
	}

	source := "neva"
	severity := protocol.DiagnosticSeverityError

	return protocol.PublishDiagnosticsParams{
		URI: uri,
		Diagnostics: []protocol.Diagnostic{
			{
				Range:    startStopRange,
				Severity: &severity,
				Source:   &source,
				Message:  compilerErr.Error(),
				Data:     time.Now(),
			},
		},
	}
}
