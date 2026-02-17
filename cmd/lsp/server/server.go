package server

import (
	"context"
	"path/filepath"
	"sync"
	"time"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/indexer"
)

//nolint:govet // fieldalignment: preserve layout for readability.
type Server struct {
	workspacePath string
	name, version string

	handler *Handler
	logger  commonlog.Logger
	indexer indexer.Indexer

	indexMutex *sync.Mutex
	index      *src.Build

	problemsMutex *sync.Mutex
	problemFiles  map[string]struct{} // we only need to store file urls but not their problems

	activeFile      string
	activeFileMutex *sync.Mutex
}

// getBuild returns the latest indexed build snapshot when available.
// Read locking is required because indexing updates s.index concurrently with LSP request handling.
func (s *Server) getBuild() (*src.Build, bool) {
	s.indexMutex.Lock()
	defer s.indexMutex.Unlock()
	if s.index == nil {
		return nil, false
	}
	return s.index, true
}

// setBuild replaces the indexed build snapshot.
// Writers share the same lock with getBuild to avoid data races on the build pointer.
func (s *Server) setBuild(build src.Build) {
	s.indexMutex.Lock()
	s.index = &build
	s.indexMutex.Unlock()
}

// indexAndNotifyProblems does full scan of the workspace
// and sends diagnostics if there are any problems
func (s *Server) indexAndNotifyProblems(notify glsp.NotifyFunc) error {
	build, found, compilerErr := s.indexer.FullScan(
		context.Background(),
		s.workspacePath,
	)
	if !found {
		return nil
	}

	s.setBuild(build)

	if compilerErr == nil {
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
	uri := s.workspacePath
	if compilerErr.Meta != nil {
		uri = filepath.Join(s.workspacePath, compilerErr.Meta.Location.String())
	}
	s.problemFiles[uri] = struct{}{}
	notify(
		protocol.ServerTextDocumentPublishDiagnostics,
		s.createDiagnostics(*compilerErr, uri),
	)
	s.logger.Info("diagnostic sent:", "err", compilerErr)
	s.problemsMutex.Unlock()

	return nil
}

func (s *Server) createDiagnostics(
	indexerErr indexer.Error,
	uri string,
) protocol.PublishDiagnosticsParams {
	var startStopRange protocol.Range
	if indexerErr.Meta != nil {
		// If stop is 0 0, set it to the same as start but with character incremented by 1
		start := indexerErr.Meta.Start
		stop := indexerErr.Meta.Stop
		if stop.Line == 0 && stop.Column == 0 {
			stop = start
			stop.Column++
		}

		startStopRange = protocol.Range{
			Start: protocol.Position{
				Line:      toUint32(start.Line),
				Character: toUint32(start.Column),
			},
			End: protocol.Position{
				Line:      toUint32(stop.Line),
				Character: toUint32(stop.Column),
			},
		}

		// Adjust for 0-based indexing
		startStopRange.Start.Line--
		startStopRange.End.Line--
	}

	return protocol.PublishDiagnosticsParams{
		URI: uri, // uri must be full path to the file, make sure all compiler errors include full location
		Diagnostics: []protocol.Diagnostic{
			{
				Range:    startStopRange,
				Severity: ptr(protocol.DiagnosticSeverityError),
				Source:   ptr("compiler"),
				Message:  indexerErr.Message, // we don't use Error() because it will duplicate location
				Data:     time.Now(),
			},
		},
	}
}

func ptr[T any](v T) *T {
	return &v
}

func toUint32(value int) uint32 {
	if value < 0 {
		return 0
	}
	if uint64(value) > uint64(^uint32(0)) {
		return ^uint32(0)
	}
	// #nosec G115 -- bounds checked above
	return uint32(value)
}
