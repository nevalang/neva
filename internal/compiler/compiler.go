package compiler

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/nevalang/neva/pkg/ir"
)

type Compiler struct {
	repo         Repository
	parser       Parser
	desugarer    Desugarer
	analyzer     Analyzer
	srcOptimizer SrcOptimizer
	irgen        IRGenerator
	irOptimizer  IROptimizer
}

type (
	Repository interface {
		ByPath(context.Context, string) (map[string]RawPackage, error)
		Save(context.Context, string, *ir.Program) error
	}
	Parser interface {
		ParseProg(context.Context, map[string]RawPackage) (src.Program, error)
	}
	Desugarer interface {
		Desugar(src.Program) (src.Program, error)
	}
	Analyzer interface {
		Analyze(src.Program) (src.Program, error)
	}
	SrcOptimizer interface {
		Optimize(src.Program) (src.Program, error)
	}
	IRGenerator interface {
		Generate(context.Context, src.Program) (*ir.Program, error)
	}
	IROptimizer interface {
		Optimize(*ir.Program) (*ir.Program, error)
	}

	RawPackage map[string][]byte // File -> content.
)

func (c Compiler) Compile(ctx context.Context, srcPath, dstPath string) (*ir.Program, error) {
	rawProg, err := c.repo.ByPath(ctx, srcPath)
	if err != nil {
		return nil, fmt.Errorf("repo by path: %w", err)
	}

	parsedProg, err := c.parser.ParseProg(ctx, rawProg)
	if err != nil {
		return nil, fmt.Errorf("parse prog: %w", err)
	}

	desugaredProg, err := c.desugarer.Desugar(parsedProg)
	if err != nil {
		return nil, fmt.Errorf("desugar: %w", err)
	}

	analyzedProg, err := c.analyzer.Analyze(desugaredProg)
	if err != nil {
		return nil, fmt.Errorf("analyze: %w", err)
	}

	irProg, err := c.irgen.Generate(ctx, analyzedProg)
	if err != nil {
		return nil, fmt.Errorf("generate IR: %w", err)
	}

	optimizedIR, err := c.irOptimizer.Optimize(irProg)
	if err != nil {
		return nil, fmt.Errorf("optimize IR: %w", err)
	}

	if dstPath == "" {
		return irProg, nil
	}

	if err := c.repo.Save(ctx, dstPath, optimizedIR); err != nil {
		return nil, fmt.Errorf("repo save: %w", err)
	}

	return irProg, nil
}

func New(
	repo Repository,
	parser Parser,
	desugarer Desugarer,
	analyzer Analyzer,
	srcOptimizer SrcOptimizer,
	irgen IRGenerator,
	irOptimizer IROptimizer,
) Compiler {
	return Compiler{
		repo:         repo,
		parser:       parser,
		desugarer:    desugarer,
		analyzer:     analyzer,
		srcOptimizer: srcOptimizer,
		irgen:        irgen,
		irOptimizer:  irOptimizer,
	}
}
