package server

import (
	"context"
	"path/filepath"
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

	indexMutex *sync.Mutex
	index      *src.Build

	problemsMutex *sync.Mutex
	problemFiles  map[string]struct{}

	activeFile      string
	activeFileMutex *sync.Mutex
}

func (s *Server) indexAndNotifyProblems(notify glsp.NotifyFunc) error {
	build, err := s.indexer.FullIndex(context.Background(), s.workspacePath)

	s.indexMutex.Lock()
	s.index = &build
	s.indexMutex.Unlock()

	if err == nil {
		// clear problems
		s.problemsMutex.Lock()
		for uri := range s.problemFiles {
			notify(
				protocol.ServerTextDocumentPublishDiagnostics,
				protocol.PublishDiagnosticsParams{
					URI:         uri,
					Diagnostics: []protocol.Diagnostic{},
				},
			)
		}
		s.problemFiles = make(map[string]struct{})
		s.logger.Info("full index without problems, sent empty diagnostics")
		s.problemsMutex.Unlock()
		return nil
	}

	// remember problem and send diagnostic
	s.problemsMutex.Lock()
	uri := filepath.Join(s.workspacePath, err.Meta.Location.String())
	s.problemFiles[uri] = struct{}{}
	notify(
		protocol.ServerTextDocumentPublishDiagnostics,
		s.createDiagnostics(*err, uri),
	)
	s.logger.Info("diagnostic sent:", "err", err)
	s.problemsMutex.Unlock()

	return nil
}

func (s *Server) createDiagnostics(compilerErr compiler.Error, uri string) protocol.PublishDiagnosticsParams {
	var startStopRange protocol.Range
	if compilerErr.Meta != nil {
		// If stop is 0 0, set it to the same as start but with character incremented by 1
		if compilerErr.Meta.Stop.Line == 0 && compilerErr.Meta.Stop.Column == 0 {
			compilerErr.Meta.Stop = compilerErr.Meta.Start
			compilerErr.Meta.Stop.Column++
		}

		startStopRange = protocol.Range{
			Start: protocol.Position{
				Line:      uint32(compilerErr.Meta.Start.Line),
				Character: uint32(compilerErr.Meta.Start.Column),
			},
			End: protocol.Position{
				Line:      uint32(compilerErr.Meta.Stop.Line),
				Character: uint32(compilerErr.Meta.Stop.Column),
			},
		}

		// Adjust for 0-based indexing
		startStopRange.Start.Line--
		startStopRange.End.Line--
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
