package interpreter

import (
	"context"

	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/pkg/ir"
)

type Interpreter struct {
	compiler Compiler
	adapter  Adapter
	runtime  Runtime
}

type (
	Compiler interface {
		Compile(ctx context.Context, src, dst string) (*ir.Program, error)
	}
	Adapter interface {
		Adapt(irProg *ir.Program) (runtime.Program, error)
	}
	Runtime interface {
		Run(context.Context, runtime.Program) (code int, err error)
	}
)

func (i Interpreter) Interpret(ctx context.Context, path string) (int, error) {
	irProg, err := i.compiler.Compile(ctx, path, "")
	if err != nil {
		return 0, err
	}

	rprog, err := i.adapter.Adapt(irProg)
	if err != nil {
		return 0, err
	}

	code, err := i.runtime.Run(ctx, rprog)
	if err != nil {
		return 0, err
	}

	return code, nil
}

func New(
	compiler Compiler,
	adapter Adapter,
	runtime Runtime,
) Interpreter {
	return Interpreter{
		compiler: compiler,
		adapter:  adapter,
		runtime:  runtime,
	}
}
