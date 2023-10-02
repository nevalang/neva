package compiler

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/nevalang/neva/pkg/ir"
)

type Compiler struct {
	repo     Repository
	parser   Parser
	analyzer Analyzer
	irgen    IRGenerator
}

type (
	Repository interface {
		ByPath(context.Context, string) (map[string]RawPackage, error)
		Save(context.Context, string, *ir.Program) error
	}
	Parser interface {
		ParseFiles(context.Context, map[string][]byte) (map[string]src.File, error)
	}
	Analyzer interface {
		Analyze(src.Program) error
	}
	IRGenerator interface {
		Generate(context.Context, src.Package) (*ir.Program, error)
	}

	RawPackage map[string][]byte // File -> content.
)

func (c Compiler) Compile(ctx context.Context, srcPath, dstPath string) (*ir.Program, error) {
	raw, err := c.repo.ByPath(ctx, srcPath)
	if err != nil {
		return nil, fmt.Errorf("by path: %w", err)
	}

	parsedPackages := make(src.Program, len(raw))
	for pkgName, files := range raw {
		parsedFiles, err := c.parser.ParseFiles(ctx, files)
		if err != nil {
			return nil, fmt.Errorf("parse files: %w", err)
		}
		parsedPackages[pkgName] = parsedFiles
	}

	if err := c.analyzer.Analyze(parsedPackages); err != nil {
		return nil, fmt.Errorf("analyze: %w", err)
	}

	irProg, err := c.irgen.Generate(ctx, parsedPackages["main"]) // TODO use all packages
	if err != nil {
		return nil, fmt.Errorf("generate: %w", err)
	}

	if dstPath == "" {
		return irProg, nil
	}

	if err := c.repo.Save(ctx, dstPath, irProg); err != nil {
		return nil, fmt.Errorf("repo save: %w", err)
	}

	return irProg, nil
}

func New(
	repo Repository,
	parser Parser,
	analyzer Analyzer,
	irgen IRGenerator,
) Compiler {
	return Compiler{
		repo:     repo,
		parser:   parser,
		analyzer: analyzer,
		irgen:    irgen,
	}
}
