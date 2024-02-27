package interpreter

import (
	"context"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/runtime/adapter"
	"github.com/nevalang/neva/internal/runtime/funcs"
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
	builder builder.Builder,
	compiler compiler.Compiler,
	isDebug bool,
) Interpreter {
	var connector runtime.Connector
	if isDebug {
		connector = runtime.NewConnector(DebugEventListener{})
	} else {
		connector = runtime.NewDefaultConnector()
	}
	return Interpreter{
		builder:  builder,
		compiler: compiler,
		adapter:  adapter.NewAdapter(),
		runtime: runtime.New(
			connector,
			runtime.MustNewFuncRunner(
				funcs.CreatorRegistry(),
			),
		),
	}
}
