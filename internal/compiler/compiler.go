package compiler

import (
	"context"
	"errors"

	"github.com/nevalang/neva/internal/src"
	"github.com/nevalang/neva/pkg/ir"
)

type Compiler struct {
	analyzer Analyzer
	irgen    IRGenerator
}

type (
	Analyzer interface {
		Analyze(context.Context, map[string]src.File) (map[string]src.File, error)
	}
	IRGenerator interface {
		Generate(context.Context, map[string]src.File) (*ir.Program, error)
	}
)

var (
	ErrAnalyzer = errors.New("analyzer")
	ErrIrGen    = errors.New("ir generator")
	ErrBackend  = errors.New("backend")
)

func (c Compiler) Compile(ctx context.Context, prog map[string]src.File) (*ir.Program, error) {
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
