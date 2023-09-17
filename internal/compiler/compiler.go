package compiler

import (
	"context"
	"errors"

	ir "github.com/nevalang/neva/pkg/ir/api"
)

type Compiler struct {
	analyzer Analyzer
	irgen    IRGenerator
}

type (
	Analyzer interface {
		Analyze(context.Context, HighLvlProgram) (HighLvlProgram, error)
	}
	IRGenerator interface {
		Generate(context.Context, HighLvlProgram) (*ir.Program, error)
	}
)

var (
	ErrAnalyzer = errors.New("analyzer")
	ErrIrGen    = errors.New("ir generator")
	ErrBackend  = errors.New("backend")
)

func (c Compiler) Compile(ctx context.Context, prog HighLvlProgram) (*ir.Program, error) {
	analyzedProg, err := c.analyzer.Analyze(ctx, prog)
	if err != nil {
		return nil, errors.Join(ErrAnalyzer, err)
	}

	irProg, err := c.irgen.Generate(ctx, analyzedProg)
	if err != nil {
		return nil, errors.Join(ErrIrGen, err)
	}

	return irProg, nil
}
