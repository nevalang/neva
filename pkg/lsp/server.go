package lsp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/nevalang/neva/internal/compiler/analyzer"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/pkg/lsp/indexer"
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
func (s *Server) saveIndex(build src.Build) {
	s.mu.Lock()
	s.index = &build
	s.mu.Unlock()
}

func (s *Server) indexAndNotifyProblems(notify glsp.NotifyFunc) error {
	build, analyzerErr, err := s.indexer.FullIndex(context.Background(), s.workspacePath)
	if err != nil {
		return fmt.Errorf("%w: index", err)
	}
	s.saveIndex(build)

	if analyzerErr == nil {
		notify(
			protocol.ServerTextDocumentPublishDiagnostics,
			protocol.PublishDiagnosticsParams{}, // clear problems
		)
		s.logger.Info("full index without problems, sent empty diagnostics")
		return nil
	}

	notify(
		protocol.ServerTextDocumentPublishDiagnostics,
		s.createDiagnostics(*analyzerErr),
	)

	s.logger.Info("diagnostic sent: " + analyzerErr.Error())

	return nil
}

func (s *Server) createDiagnostics(analyzerErr analyzer.Error) protocol.PublishDiagnosticsParams {
	source := "neva"
	severity := protocol.DiagnosticSeverityError

	var uri string
	if analyzerErr.Location != nil {
		uri = fmt.Sprintf(
			"%s/%s/%s",
			s.workspacePath,
			analyzerErr.Location.PkgName,
			analyzerErr.Location.FileName+".neva",
		)
	}

	var protocolRange protocol.Range
	if analyzerErr.Meta != nil {
		protocolRange = protocol.Range{
			Start: protocol.Position{
				Line:      uint32(analyzerErr.Meta.Start.Line),
				Character: uint32(analyzerErr.Meta.Start.Column),
			},
			End: protocol.Position{
				Line:      uint32(analyzerErr.Meta.Stop.Line),
				Character: uint32(analyzerErr.Meta.Stop.Column),
			},
		}
	}

	return protocol.PublishDiagnosticsParams{
		URI: uri,
		Diagnostics: []protocol.Diagnostic{
			{
				Range:    protocolRange,
				Severity: &severity,
				Source:   &source,
				Message:  analyzerErr.Error(),
				Data:     time.Now(),
				// Unused:
				Tags:               []protocol.DiagnosticTag{},
				Code:               &protocol.IntegerOrString{Value: nil},
				CodeDescription:    &protocol.CodeDescription{HRef: ""},
				RelatedInformation: []protocol.DiagnosticRelatedInformation{},
			},
		},
	}
}
