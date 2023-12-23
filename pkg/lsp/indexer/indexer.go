package indexer

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/desugarer"
	"github.com/nevalang/neva/internal/compiler/parser"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

type Indexer struct {
	builder   builder.Builder
	parser    parser.Parser
	desugarer desugarer.Desugarer
	analyzer  analyzer.Analyzer
}

func (i Indexer) FullIndex(ctx context.Context, path string) (src.Build, *analyzer.Error, error) {
	rawBuild, err := i.builder.Build(ctx, path)
	if err != nil {
		return src.Build{}, nil, fmt.Errorf("builder: %w", err)
	}

	parsedMods, err := i.parser.ParseModules(rawBuild.Modules)
	if err != nil {
		return src.Build{}, nil, fmt.Errorf("parse prog: %w", err)
	}

	parsedBuild := src.Build{
		EntryModRef: rawBuild.EntryModRef,
		Modules:     parsedMods,
	}

	_, err = i.analyzer.AnalyzeBuild(parsedBuild)
	if err == nil {
		return parsedBuild, nil, nil
	}

	analyzerErr, ok := err.(*analyzer.Error)
	if !ok {
		return src.Build{}, nil, fmt.Errorf("cast analyzer err: %w", err)
	}

	return parsedBuild, analyzerErr, nil
}

func New(
	builder builder.Builder,
	parser parser.Parser,
	desugarer desugarer.Desugarer,
	analyzer analyzer.Analyzer,
) Indexer {
	return Indexer{
		builder:   builder,
		parser:    parser,
		desugarer: desugarer,
		analyzer:  analyzer,
	}
}
