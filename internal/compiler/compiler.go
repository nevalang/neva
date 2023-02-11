package compiler

import (
	"context"

	"github.com/emil14/neva/internal/compiler/src"
	"github.com/emil14/neva/pkg/tools"
)

type Compiler[T any] struct {
	builder     Builder
	analyzer    Analyzer
	synthesizer Synthesizer[T]
}

type (
	Builder interface {
		Build(context.Context, string) (src.Program, error)
	}
	Analyzer interface {
		Analyze(context.Context, src.Program) (src.Program, error)
	}
	Synthesizer[T any] interface {
		Synthesize(context.Context, src.Program) (T, error)
	}
)

func (c Compiler[T]) Compile(ctx context.Context, path string) (*T, error) {
	prog, err := c.builder.Build(ctx, path)
	if err != nil {
		return nil, err
	}

	prog, err = c.analyzer.Analyze(ctx, prog)
	if err != nil {
		return nil, err
	}

	rprog, err := c.synthesizer.Synthesize(ctx, prog)
	if err != nil {
		return nil, err
	}

	return &rprog, err
}

func MustNew[T any](b Builder, a Analyzer, s Synthesizer[T]) Compiler[T] {
	tools.NilPanic(b, a, s)

	return Compiler[T]{
		builder:     b,
		analyzer:    a,
		synthesizer: s,
	}
}
