package compiler

import (
	"context"
	"io"

	"github.com/emil14/neva/internal/compiler/ir"
	"github.com/emil14/neva/internal/compiler/src"
)

type Compiler struct {
	analyzer Analyzer
	irgen    IRGenerator
	backend  Backend
	writer   io.Writer
}

type (
	Analyzer interface {
		Analyze(context.Context, src.Program) (src.Program, error)
	}
	IRGenerator interface {
		Generate(context.Context, src.Program) (ir.Program, error)
	}
	Backend interface {
		GenerateTarget(context.Context, ir.Program) ([]byte, error)
	}
)

func (c Compiler) Compile(ctx context.Context, prog src.Program) ([]byte, error) {
	analyzedProg, err := c.analyzer.Analyze(ctx, prog)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	irProg, err := c.irgen.Generate(ctx, analyzedProg)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	target, err := c.backend.GenerateTarget(ctx, irProg)
	if err != nil {
		return nil, err
	}

	_, err = c.writer.Write(target)
	if err != nil {
		return nil, err
	}

	return target, nil
}
