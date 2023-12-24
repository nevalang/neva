package main

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/compiler"
	builder "github.com/nevalang/neva/internal/pkgmanager"
	"github.com/nevalang/neva/internal/runtime"
)

type Interpreter struct {
	builder  builder.PkgManager
	compiler compiler.Compiler
	runtime  runtime.Runtime
	adapter  Adapter
}

func (i Interpreter) Interpret(ctx context.Context, workdirPath string, mainPkgName string) (int, error) {
	build, err := i.builder.Build(ctx, workdirPath)
	if err != nil {
		return 0, fmt.Errorf("build: %w", err)
	}

	irProg, err := i.compiler.Compile(ctx, build, workdirPath, mainPkgName)
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
	compiler compiler.Compiler,
	adapter Adapter,
	runtime runtime.Runtime,
	builder builder.PkgManager,
) Interpreter {
	return Interpreter{
		compiler: compiler,
		adapter:  adapter,
		runtime:  runtime,
		builder:  builder,
	}
}
