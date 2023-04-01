package compiler

import (
	"context"
	"errors"

	"github.com/nevalang/nevalang/internal/compiler/ir"
)

type Compiler struct {
	analyzer Analyzer
	irgen    IRGenerator
	backend  Backend
}

type (
	Analyzer interface {
		Analyze(context.Context, Program) (Program, error)
	}
	IRGenerator interface {
		Generate(context.Context, Program) (ir.Program, error)
	}
	Backend interface {
		GenerateTarget(context.Context, ir.Program) ([]byte, error)
	}
)

var (
	ErrAnalyzer = errors.New("analyzer")
	ErrIrGen    = errors.New("ir generator")
	ErrBackend  = errors.New("backend")
)

func (c Compiler) Compile(ctx context.Context, prog Program) ([]byte, error) {
	analyzedProg, err := c.analyzer.Analyze(ctx, prog)
	if err != nil {
		return nil, errors.Join(ErrAnalyzer, err)
	}

	irProg, err := c.irgen.Generate(ctx, analyzedProg)
	if err != nil {
		return nil, errors.Join(ErrIrGen, err)
	}

	target, err := c.backend.GenerateTarget(ctx, irProg)
	if err != nil {
		return nil, errors.Join(ErrBackend, err)
	}

	return target, nil
}
