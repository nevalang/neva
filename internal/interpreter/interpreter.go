package interpreter

import (
	"context"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/interpreter/adapter"
	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/funcs"
)

type Interpreter struct {
	builder  builder.Builder
	compiler compiler.Compiler
	runtime  runtime.Runtime
	adapter  adapter.Adapter
}

func (i Interpreter) Interpret(ctx context.Context, main string, debug bool) *compiler.Error {
	irProg, compilerErr := i.compiler.CompileToIR(main, debug)
	if compilerErr != nil {
		return compiler.Error{
			Location: &sourcecode.Location{
				PkgName: main,
			},
		}.Wrap(compilerErr)
	}

	rprog, err := i.adapter.Adapt(irProg)
	if err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &sourcecode.Location{PkgName: main},
		}
	}

	if err := i.runtime.Run(ctx, rprog); err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &sourcecode.Location{PkgName: main},
		}
	}

	return nil
}

func New(
	builder builder.Builder,
	compiler compiler.Compiler,
) Interpreter {
	return Interpreter{
		builder:  builder,
		compiler: compiler,
		adapter:  adapter.NewAdapter(),
		runtime: runtime.New(
			runtime.NewFuncRunner(
				funcs.CreatorRegistry(),
			),
		),
	}
}
