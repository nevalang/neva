package indexer

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/parser"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/tliron/commonlog"
)

type Indexer struct {
	Builder  builder.Builder
	Parser   parser.Parser
	Analyzer analyzer.Analyzer
	Logger   commonlog.Logger
}

func (i Indexer) FullIndex(ctx context.Context, path string) (mod src.Module, analyzerErr *analyzer.Error, err error) {
	build, err := i.Builder.Build(ctx, path)
	if err != nil {
		return src.Module{}, nil, fmt.Errorf("builder: %w", err)
	}

	rawMod := build.Modules[build.EntryModule] // TODO use all mods

	parsedPkgs, err := i.Parser.ParsePackages(ctx, rawMod.Packages)
	if err != nil {
		return src.Module{}, nil, fmt.Errorf("parse prog: %w", err)
	}

	mod = src.Module{
		Manifest: rawMod.Manifest,
		Packages: parsedPkgs,
	}

	// we interpret analyzer error as a message, not failure
	if _, err = i.Analyzer.Analyze(mod); err != nil {
		analyzerErr, ok := err.(*analyzer.Error) // FIXME for some reason we loose info after cast
		if !ok {
			i.Logger.Errorf("Analyzer returned an error of unexpected type: %T", err)
		}
		return mod, analyzerErr, nil
	}

	return mod, nil, nil
}
