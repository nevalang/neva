package lsp

import (
	"context"
	"fmt"
	"sync"
	"time"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/pkg/lsp/indexer"
	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Server struct {
	rootPath      string
	name, version string

	handler *Handler
	logger  commonlog.Logger
	indexer indexer.Indexer

	mu    *sync.Mutex
	state *State
}

type State struct {
	mod      src.Module
	problems string
}

// setState allows to update state in a thread-safe manner.
func (s *Server) setState(mod src.Module, problem string) {
	s.mu.Lock()
	s.state = &State{
		mod:      mod,
		problems: problem,
	}
	s.mu.Unlock()
}

func (s *Server) indexAndNotifyProblems(notify glsp.NotifyFunc) error {
	prog, problems, err := s.indexer.FullIndex(context.Background(), s.rootPath)
	if err != nil {
		return fmt.Errorf("%w: index", err)
	}
	s.setState(prog, problems)

	if s.state.problems == "" {
		notify(
			protocol.ServerTextDocumentPublishDiagnostics,
			protocol.PublishDiagnosticsParams{}, // clear problems
		)
		s.logger.Info("full index without problems, sent empty diagnostics")
		return nil
	}

	source := "neva"
	severity := protocol.DiagnosticSeverityError

	notify(
		protocol.ServerTextDocumentPublishDiagnostics,
		protocol.PublishDiagnosticsParams{
			URI: "",
			Diagnostics: []protocol.Diagnostic{
				{
					Range: protocol.Range{
						Start: protocol.Position{
							Line:      0,
							Character: 0,
						},
						End: protocol.Position{
							Line:      0,
							Character: 0,
						},
					},
					Severity: &severity,
					Code: &protocol.IntegerOrString{
						Value: nil,
					},
					CodeDescription: &protocol.CodeDescription{
						HRef: "",
					},
					Source:             &source,
					Message:            s.state.problems,
					Tags:               []protocol.DiagnosticTag{},
					RelatedInformation: []protocol.DiagnosticRelatedInformation{},
					Data:               time.Now(),
				},
			},
		},
	)

	s.logger.Info("diagnostic sent: " + s.state.problems)

	return nil
}
