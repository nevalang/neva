package indexer

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/tliron/commonlog"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/analyzer"
	"github.com/nevalang/neva/internal/compiler/parser"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
)

type Indexer struct {
	fe       compiler.Frontend
	analyzer analyzer.Analyzer
	logger   commonlog.Logger
}

func (i Indexer) FullScan(
	ctx context.Context,
	workspacePath string,
) (src.Build, bool, *compiler.Error) {
	feResult, err := i.fe.Process(ctx, workspacePath)
	if err != nil {
		return src.Build{}, false, err
	}

	if isParentPath(workspacePath, feResult.Path) {
		i.logger.Debug(
			"nevalang module found but not part of workspace",
			"path", feResult.Path, "workspacePath", workspacePath,
		)
		return src.Build{}, false, nil
	}

	i.logger.Debug("nevalang module found in workspace", "path", feResult.Path)

	aBuild, err := i.analyzer.AnalyzeBuild(feResult.ParsedBuild)
	if err != nil {
		return src.Build{}, true, err.Unwrap() // use only deepest compiler error for now
	}

	return aBuild, true, nil
}

func isParentPath(parent, child string) bool {
	parent = filepath.Clean(parent)
	child = filepath.Clean(child)

	rel, err := filepath.Rel(parent, child)
	if err != nil {
		panic(err)
	}
	if rel == "." {
		return false
	}

	return !strings.HasPrefix(rel, "..")
}

func New(
	builder builder.Builder,
	parser parser.Parser,
	analyzer analyzer.Analyzer,
	logger commonlog.Logger,
) Indexer {
	return Indexer{
		fe:       compiler.NewFrontend(builder, parser),
		analyzer: analyzer,
		logger:   logger,
	}
}
