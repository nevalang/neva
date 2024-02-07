package interpreter

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	builder "github.com/nevalang/neva/internal/pkgmanager"
	"github.com/nevalang/neva/internal/runtime"
)

type Interpreter struct {
	builder  builder.Manager
	compiler compiler.Compiler
	runtime  runtime.Runtime
	adapter  Adapter
}

func (i Interpreter) Interpret(ctx context.Context, workdirPath string, mainPkgName string) error {
	build, err := i.builder.Build(ctx, workdirPath)
	if err != nil {
		return fmt.Errorf("build: %w", err)
	}

	irProg, compilerErr := i.compiler.CompileToIR(ctx, build, workdirPath, mainPkgName)
	if compilerErr != nil {
		return compilerErr
	}

	rprog, err := i.adapter.Adapt(irProg)
	if err != nil {
		return err
	}

	if err := i.runtime.Run(ctx, rprog); err != nil {
		return err
	}

	return nil
}

func New(
	compiler compiler.Compiler,
	adapter Adapter,
	runtime runtime.Runtime,
	builder builder.Manager,
) Interpreter {
	return Interpreter{
		compiler: compiler,
		adapter:  adapter,
		runtime:  runtime,
		builder:  builder,
	}
}
