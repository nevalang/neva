package interpreter

import (
	"context"

	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/pkgmanager"
	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/pkg/sourcecode"
)

type Interpreter struct {
	builder  pkgmanager.Manager
	compiler compiler.Compiler
	runtime  runtime.Runtime
	adapter  Adapter
}

func (i Interpreter) Interpret(ctx context.Context, workdirPath string, mainPkgName string) *compiler.Error {
	irProg, compilerErr := i.compiler.CompileToIR(workdirPath, mainPkgName)
	if compilerErr != nil {
		return compiler.Error{
			Location: &sourcecode.Location{
				PkgName: mainPkgName,
			},
		}.Merge(compilerErr)
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
	adapter Adapter,
	runtime runtime.Runtime,
	builder pkgmanager.Manager,
) Interpreter {
	return Interpreter{
		compiler: compiler,
		adapter:  adapter,
		runtime:  runtime,
		builder:  builder,
	}
}
