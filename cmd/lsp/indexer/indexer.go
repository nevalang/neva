package indexer

import (
	"context"
	"path/filepath"
	"strings"

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

func (i Indexer) FullScan(
	ctx context.Context,
	workspacePath string,
) (src.Build, bool, *compiler.Error) {
	feResult, err := i.fe.Process(ctx, workspacePath)
	if err != nil {
		return src.Build{}, false, err
	}

	// if nevalang module is found, but it's not part of the workspace
	if isParentPath(workspacePath, feResult.Path) {
		return src.Build{}, false, nil
	}

	aBuild, err := i.analyzer.AnalyzeBuild(feResult.ParsedBuild)
	if err != nil {
		return src.Build{}, false, err
	}

	return aBuild, true, nil
}

func isParentPath(parent, child string) bool {
	parent = filepath.Clean(parent)
	child = filepath.Clean(child)

	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}

	return !strings.HasPrefix(rel, "..")
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
