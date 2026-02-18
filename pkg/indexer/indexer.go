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
	src "github.com/nevalang/neva/pkg/ast"
)

const mainPackageNotFoundMessage = "main package not found"

// Indexer performs workspace scans and returns analyzed source build snapshots.
type Indexer struct {
	fe       compiler.Frontend
	analyzer analyzer.Analyzer
	logger   commonlog.Logger
}

// FullScan processes and analyzes a Neva workspace.
func (i Indexer) FullScan(
	ctx context.Context,
	workspacePath string,
) (src.Build, bool, *Error) {
	feResult, err := i.fe.Process(ctx, workspacePath)
	if err != nil {
		return src.Build{}, false, wrapCompilerError(err)
	}

	if isParentPath(workspacePath, feResult.Path) {
		i.logger.Debug(
			"nevalang module found but not part of workspace",
			"path", feResult.Path, "workspacePath", workspacePath,
		)
		return src.Build{}, false, nil
	}

	i.logger.Debug("nevalang module found in workspace", "path", feResult.Path)

	aBuild, err := i.analyzer.Analyze(feResult.ParsedBuild, feResult.MainPkg)
	if err != nil {
		if !isMainPackageNotFoundError(err) {
			return src.Build{}, true, wrapCompilerError(err)
		}

		// Workspace indexing should remain useful even when the workspace root
		// is not itself a runnable entry package.
		i.logger.Info(
			"main package not found; falling back to workspace library analysis",
			"mainPkg", feResult.MainPkg,
			"workspacePath", workspacePath,
		)

		fallbackBuild, fallbackErr := i.analyzer.Analyze(feResult.ParsedBuild, "")
		if fallbackErr != nil {
			i.logger.Warning(
				"workspace fallback analysis failed",
				"mainPkg", feResult.MainPkg,
				"workspacePath", workspacePath,
				"err", fallbackErr,
			)
			return src.Build{}, true, wrapCompilerError(fallbackErr)
		}

		i.logger.Info(
			"workspace fallback analysis succeeded",
			"mainPkg", feResult.MainPkg,
			"workspacePath", workspacePath,
		)
		return fallbackBuild, true, nil
	}

	return aBuild, true, nil
}

func isMainPackageNotFoundError(err *compiler.Error) bool {
	if err == nil {
		return false
	}
	return err.Unwrap().Message == mainPackageNotFoundMessage
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

// New constructs Indexer from compiler frontend dependencies.
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
