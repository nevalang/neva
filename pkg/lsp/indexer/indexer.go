package indexer

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/desugarer"
	"github.com/nevalang/neva/internal/compiler/parser"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/tliron/commonlog"
)

type Indexer struct {
	builder   builder.Builder
	parser    parser.Parser
	desugarer desugarer.Desugarer
	analyzer  analyzer.Analyzer
	logger    commonlog.Logger
}

func (i Indexer) FullIndex(ctx context.Context, path string) (mod src.Module, analyzerErr *analyzer.Error, err error) {
	build, err := i.builder.Build(ctx, path)
	if err != nil {
		return src.Module{}, nil, fmt.Errorf("builder: %w", err)
	}

	rawMod := build.Modules[build.EntryModRef]

	parsedPkgs, err := i.parser.ParsePackages(ctx, rawMod.Packages)
	if err != nil {
		return src.Module{}, nil, fmt.Errorf("parse prog: %w", err)
	}

	mod = src.Module{
		Manifest: rawMod.Manifest,
		Packages: parsedPkgs,
	}

	// we interpret analyzer error as a message, not failure
	if _, err = i.analyzer.Analyze(mod); err != nil {
		analyzerErr, ok := err.(*analyzer.Error) // FIXME for some reason we loose info after cast
		if !ok {
			i.logger.Errorf("Analyzer returned an error of unexpected type: %T", err)
		}
		return mod, analyzerErr, nil
	}

	return mod, nil, nil
}

func New(
	builder builder.Builder,
	parser parser.Parser,
	desugarer desugarer.Desugarer,
	analyzer analyzer.Analyzer,
	logger commonlog.Logger,
) Indexer {
	return Indexer{
		builder:   builder,
		parser:    parser,
		desugarer: desugarer,
		analyzer:  analyzer,
		logger:    logger,
	}
}
