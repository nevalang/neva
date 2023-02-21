package compiler

import (
	"context"

	"github.com/emil14/neva/internal/compiler/src"
)

type Compiler[T any] struct {
	analyzer Analyzer
	// optimizer   Optimizer
	synthesizer Synthesizer[T]
}

type (
	Analyzer interface {
		Analyze(context.Context, src.Prog) (src.Prog, error)
	}
	Optimizer interface {
		Optimize(context.Context, src.Prog) (src.Prog, error)
	}
	Synthesizer[T any] interface {
		Synthesize(context.Context, src.Prog) (T, error)
	}
)

func (c Compiler[T]) Compile(ctx context.Context, resolvedProg src.Prog) (*T, error) {
	resolvedProg, err := c.analyzer.Analyze(ctx, resolvedProg)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	target, err := c.synthesizer.Synthesize(ctx, resolvedProg)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return &target, nil
}

// func MustNew[T any](b Builder, a Analyzer, s Synthesizer[T]) Compiler[T] {
// 	tools.NilPanic(b, a, s)

// 	return Compiler[T]{
// 		builder:     b,
// 		analyzer:    a,
// 		synthesizer: s,
// 	}
// }
