package compiler

import (
	"context"
	"errors"
)

type Compiler struct {
	analyzer Analyzer
	llrgen   LowLvlGenerator
	backend  Backend
}

type (
	Analyzer interface {
		Analyze(context.Context, HighLvlProgram) (HighLvlProgram, error)
	}
	LowLvlGenerator interface {
		Generate(context.Context, HighLvlProgram) (LowLvlProgram, error)
	}
	Backend interface {
		GenerateTarget(context.Context, LowLvlProgram) error
	}
)

var (
	ErrAnalyzer = errors.New("analyzer")
	ErrIrGen    = errors.New("ir generator")
	ErrBackend  = errors.New("backend")
)

func (c Compiler) Compile(ctx context.Context, prog HighLvlProgram) error {
	analyzedProg, err := c.analyzer.Analyze(ctx, prog)
	if err != nil {
		return errors.Join(ErrAnalyzer, err)
	}

	irProg, err := c.llrgen.Generate(ctx, analyzedProg)
	if err != nil {
		return errors.Join(ErrIrGen, err)
	}

	if err := c.backend.GenerateTarget(ctx, irProg); err != nil {
		return errors.Join(ErrBackend, err)
	}

	return nil
}
