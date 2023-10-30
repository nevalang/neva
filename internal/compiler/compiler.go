package compiler

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/compiler/src"
	"github.com/nevalang/neva/pkg/ir"
)

type Compiler struct {
	front FrontEnd
	back  IRGen
}

type FrontEnd struct {
	parser   Parser
	analyzer Analyzer
}

func NewFrontEnd(parser Parser, analyzer Analyzer) FrontEnd {
	return FrontEnd{
		parser:   parser,
		analyzer: analyzer,
	}
}

// mainPkgName can be "", in that case no executable package is required
func (f FrontEnd) Process(ctx context.Context, rawProg map[string]RawPackage, mainPkgName string) (src.Program, error) {
	parsedProg, err := f.parser.Parse(ctx, rawProg)
	if err != nil {
		return nil, fmt.Errorf("parse prog: %w", err)
	}

	var analyzedProg src.Program
	if mainPkgName == "" {
		analyzedProg, err = f.analyzer.Analyze(parsedProg)
		if err != nil {
			return nil, fmt.Errorf("analyze: %w", err)
		}
	} else {
		analyzedProg, err = f.analyzer.AnalyzeExecutable(parsedProg, mainPkgName)
		if err != nil {
			return nil, fmt.Errorf("analyze: %w", err)
		}
	}

	return analyzedProg, nil
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
	prog, err := c.front.Process(ctx, rawProg, mainPkgName)
	if err != nil {
		return nil, fmt.Errorf("frontend: %w", err)
	}

	irProg, err := c.back.Generate(ctx, prog)
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
		front: FrontEnd{
			parser:   parser,
			analyzer: analyzer,
		},
		back: irgen,
	}
}
