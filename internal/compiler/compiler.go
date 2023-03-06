package compiler

import (
	"context"

	"github.com/emil14/neva/internal/compiler/ir"
	"github.com/emil14/neva/internal/compiler/src"
)

type Compiler struct {
	analyzer Analyzer
	irgen    IRGenerator
	backend  Backend
}

type (
	Analyzer interface {
		Analyze(context.Context, src.Prog) (src.Prog, error) // returns program ready for synthesis
	}
	IRGenerator interface {
		GenerateIR(context.Context, src.Prog) (ir.Program, error)
	}
	Backend interface {
		GenerateTarget(context.Context, ir.Program) ([]byte, error)
	}
)

func (c Compiler) Compile(ctx context.Context, analyzedSrc src.Prog) ([]byte, error) {
	analyzedSrc, err := c.analyzer.Analyze(ctx, analyzedSrc)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	irprog, err := c.irgen.GenerateIR(ctx, analyzedSrc)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	target, err := c.backend.GenerateTarget(ctx, irprog)
	if err != nil {
		return nil, err
	}

	return target, nil
}
