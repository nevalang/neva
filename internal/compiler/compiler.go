package compiler

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/nevalang/neva/pkg/ir"
)

type Compiler struct {
	parser   Parser
	analyzer Analyzer
	irgen    IRGen
}

type (
	Parser interface {
		Parse(context.Context, map[string]RawPackage) (src.Program, error)
	}

	RawPackage map[string][]byte

	Analyzer interface {
		Analyze(prog src.Program) (src.Program, error)
		AnalyzeExecutable(prog src.Program, mainPkg string) (src.Program, error)
	}
)

type IRGen interface {
	Generate(context.Context, src.Program) (*ir.Program, error)
}

// Compile is like Analyze but also produces IR.
func (c Compiler) Compile(ctx context.Context, rawProg map[string]RawPackage, mainPkgName string) (*ir.Program, error) {
	parsedProg, err := c.parser.Parse(ctx, rawProg)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	analyzedProg, err := c.analyzer.Analyze(parsedProg)
	if err != nil {
		return nil, fmt.Errorf("analyzer: %w", err)
	}

	irProg, err := c.irgen.Generate(ctx, analyzedProg)
	if err != nil {
		return nil, fmt.Errorf("generate IR: %w", err)
	}

	return irProg, nil
}

// New creates new Compiler instance. You can omit irgen if all you need is Analyze method.
func New(
	parser Parser,
	analyzer Analyzer,
	irgen IRGen,
) Compiler {
	return Compiler{
		parser:   parser,
		analyzer: analyzer,
		irgen:    irgen,
	}
}
