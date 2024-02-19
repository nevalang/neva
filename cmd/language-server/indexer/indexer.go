package indexer

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/desugarer"
	"github.com/nevalang/neva/internal/compiler/parser"
	"github.com/nevalang/neva/internal/pkgmanager"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

type Indexer struct {
	pkgManager pkgmanager.Manager
	parser     parser.Parser
	desugarer  desugarer.Desugarer
	analyzer   analyzer.Analyzer
}

func (i Indexer) FullIndex(ctx context.Context, path string) (src.Build, *compiler.Error, error) {
	rawBuild, err := i.pkgManager.Build(ctx, path)
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

	return parsedBuild, err, nil
}

func New(
	builder pkgmanager.Manager,
	parser parser.Parser,
	desugarer desugarer.Desugarer,
	analyzer analyzer.Analyzer,
) Indexer {
	return Indexer{
		pkgManager: builder,
		parser:     parser,
		desugarer:  desugarer,
		analyzer:   analyzer,
	}
}
