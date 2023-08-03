package interpreter

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/shared"
	"github.com/nevalang/neva/pkg/tools"
)

type Interpreter struct {
	parser      Parser
	llrgen      LowLvlGenerator
	transformer Transformer
	runtime     Runtime
}

type Parser interface {
	Parse(context.Context, []byte) (map[string]shared.File, error)
}

type LowLvlGenerator interface {
	Generate(context.Context, map[string]shared.File) (shared.LowLvlProgram, error)
}

type Transformer interface {
	Transform(context.Context, shared.LowLvlProgram) (runtime.Program, error)
}

type Runtime interface {
	Run(context.Context, runtime.Program) (code int, err error)
}

func (i Interpreter) Interpret(ctx context.Context, bb []byte) (int, error) {
	hl, err := i.parser.Parse(ctx, bb)
	if err != nil {
		return 0, err
	}

	ll, err := i.llrgen.Generate(ctx, hl)
	if err != nil {
		return 0, err
	}

	rprog, err := i.transformer.Transform(ctx, ll)
	if err != nil {
		return 0, err
	}

	code, err := i.runtime.Run(ctx, rprog)
	if err != nil {
		return 0, err
	}

	return code, nil
}

func MustNew(
	parser Parser,
	llrgen LowLvlGenerator,
	transformer Transformer,
	runtime Runtime,
) Interpreter {
	tools.NilPanic(parser, llrgen, transformer, runtime)
	return Interpreter{
		parser:      parser,
		llrgen:      llrgen,
		transformer: transformer,
		runtime:     runtime,
	}
}
