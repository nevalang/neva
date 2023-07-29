package interpreter

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/shared"
)

type Interpreter struct {
	parser      Parser
	llrgen      LowLvlGenerator
	transformer Transformer
	runtime     Runtime
}

type Parser interface {
	Parse(context.Context, string) (shared.HighLvlProgram, error)
}

type LowLvlGenerator interface {
	Generate(context.Context, shared.HighLvlProgram) (shared.LowLvlProgram, error)
}

type Transformer interface {
	Transform(context.Context, shared.LowLvlProgram) (runtime.Program, error)
}

type Runtime interface {
	Run(context.Context, runtime.Program) (code int, err error)
}
