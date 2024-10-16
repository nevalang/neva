package indexer

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/parser"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

type Indexer struct {
	fe       compiler.Frontend
	analyzer analyzer.Analyzer
}

func (i Indexer) FullIndex(ctx context.Context, path string) (src.Build, *compiler.Error, error) {
	feResult, err := i.fe.Process(ctx, path)
	if err != nil {
		return src.Build{}, nil, fmt.Errorf("frontend: %w", err)
	}

	analyzedBuild, err := i.analyzer.AnalyzeBuild(feResult.ParsedBuild)
	if err != nil {
		return src.Build{}, err, nil
	}

	return analyzedBuild, nil, nil
}

func New(
	builder builder.Builder,
	parser parser.Parser,
	analyzer analyzer.Analyzer,
) Indexer {
	return Indexer{
		fe:       compiler.NewFrontend(builder, parser),
		analyzer: analyzer,
	}
}
