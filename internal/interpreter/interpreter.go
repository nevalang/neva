package interpreter

import (
	"context"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/adapter"
	"github.com/nevalang/neva/pkg/sourcecode"
)

type Interpreter struct {
	builder  builder.Builder
	compiler compiler.Compiler
	runtime  runtime.Runtime
	adapter  adapter.Adapter
}

func (i Interpreter) Interpret(ctx context.Context, workdirPath string, mainPkgName string) *compiler.Error {
	irProg, compilerErr := i.compiler.CompileToIR(workdirPath, mainPkgName)
	if compilerErr != nil {
		return compiler.Error{
			Location: &sourcecode.Location{
				PkgName: mainPkgName,
			},
		}.Wrap(compilerErr)
	}

	rprog, err := i.adapter.Adapt(irProg)
	if err != nil {
		return &compiler.Error{
			Err: err,
			Location: &sourcecode.Location{
				PkgName: mainPkgName,
			},
		}
	}

	if err := i.runtime.Run(ctx, rprog); err != nil {
		return &compiler.Error{
			Err: err,
			Location: &sourcecode.Location{
				PkgName: mainPkgName,
			},
		}
	}

	return nil
}

func New(
	compiler compiler.Compiler,
	adapter adapter.Adapter,
	runtime runtime.Runtime,
	builder builder.Builder,
) Interpreter {
	return Interpreter{
		compiler: compiler,
		adapter:  adapter,
		runtime:  runtime,
		builder:  builder,
	}
}
