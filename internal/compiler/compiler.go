package compiler

import (
	"context"
	"errors"
)

type Compiler struct {
	analyzer Analyzer
	irgen    LLRGenerator
	backend  Backend
}

type (
	Analyzer interface {
		Analyze(context.Context, HLProgram) (HLProgram, error)
	}
	LLRGenerator interface {
		Generate(context.Context, HLProgram) (LLProgram, error)
	}
	Backend interface {
		GenerateTarget(context.Context, LLProgram) ([]byte, error)
	}
)

var (
	ErrAnalyzer = errors.New("analyzer")
	ErrIrGen    = errors.New("ir generator")
	ErrBackend  = errors.New("backend")
)

func (c Compiler) Compile(ctx context.Context, prog HLProgram) ([]byte, error) {
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
