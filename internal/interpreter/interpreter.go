package interpreter

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/src"
	"github.com/nevalang/neva/pkg/ir"
)

type Interpreter struct {
	parser  SourceCodeParser
	irgen   IRGenerator
	rtgen   RuntimeProgramGenerator
	runtime Runtime
}

type SourceCodeParser interface {
	Parse(context.Context, []byte) (map[string]src.File, error)
}

type IRGenerator interface {
	Generate(context.Context, map[string]src.File) (*ir.Program, error)
}

type RuntimeProgramGenerator interface {
	Transform(context.Context, *ir.Program) (runtime.Program, error)
}

type Runtime interface {
	Run(context.Context, runtime.Program) (code int, err error)
}

func (i Interpreter) Interpret(ctx context.Context, bb []byte) (int, error) {
	// FIXME someone has to resolve imports
	hl, err := i.parser.Parse(ctx, bb)
	if err != nil {
		return 0, err
	}

	ll, err := i.irgen.Generate(ctx, hl)
	if err != nil {
		return 0, err
	}

	rprog, err := i.rtgen.Transform(ctx, ll)
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
	parser SourceCodeParser,
	irgen IRGenerator,
	transformer RuntimeProgramGenerator,
	runtime Runtime,
) Interpreter {
	if parser == nil || irgen == nil || transformer == nil || runtime == nil {
		panic("nil argument")
	}
	return Interpreter{
		parser:  parser,
		irgen:   irgen,
		rtgen:   transformer,
		runtime: runtime,
	}
}
