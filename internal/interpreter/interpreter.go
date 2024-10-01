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

// TODO remove
type Interpreter struct {
	builder  builder.Builder
	compiler compiler.Compiler
	registry map[string]runtime.FuncCreator
	adapter  adapter.Adapter
}

func (i Interpreter) Interpret(
	ctx context.Context,
	main string,
	debug bool,
	debugFile string,
) *compiler.Error {
	result, compilerErr := i.compiler.CompileToIR(main, debug)
	if compilerErr != nil {
		return compiler.Error{
			Location: &sourcecode.Location{
				PkgName: main,
			},
		}.Wrap(compilerErr)
	}

	rprog, err := i.adapter.Adapt(result.IR, debug, debugFile)
	if err != nil {
		return &compiler.Error{
			Err:      err,
			Location: &sourcecode.Location{PkgName: main},
		}
	}

	if err := runtime.Run(ctx, rprog, i.registry); err != nil {
		panic(err)
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
		registry: funcs.NewRegistry(),
	}
}
