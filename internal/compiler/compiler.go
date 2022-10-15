package compiler

import (
	"context"

	csrc "github.com/emil14/neva/internal/compiler/src"
	"github.com/emil14/neva/pkg/tools"
)

type Compiler[T any] struct {
	analyzer    Analyzer
	synthesizer Synthesizer[T]
}

type (
	Analyzer interface {
		Analyze(context.Context, csrc.Program) (csrc.Program, error)
	}
	Synthesizer[T any] interface {
		Synthesize(context.Context, csrc.Program) (T, error)
	}
)

func (c Compiler[T]) Compile(ctx context.Context, prog csrc.Program) (any, error) {
	if err := c.analyzer.Analyze(ctx, prog); err != nil {
		return nil, err
	}

	rprog, err := c.synthesizer.Synthesize(ctx, prog)
	if err != nil {
		return nil, err
	}

	return rprog, err
}

func MustNew[T any](a Analyzer, s Synthesizer[T]) Compiler[T] {
	tools.PanicWithNil(a, s)

	return Compiler[T]{
		analyzer:    a,
		synthesizer: s,
	}
}
