package interpreter

import (
	"context"
	"fmt"

	"github.com/nevalang/neva/internal/builder"
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/vm/decoder/proto"
)

type Interpreter struct {
	builder  builder.Builder
	compiler compiler.Compiler
	runtime  runtime.Runtime
	adapter  proto.Adapter
}

func (i Interpreter) Interpret(ctx context.Context, pathToMainPkg string) (int, error) {
	rawMod, err := i.builder.BuildModule(ctx, pathToMainPkg)
	if err != nil {
		return 0, fmt.Errorf("build: %w", err)
	}

	irProg, err := i.compiler.Compile(ctx, compiler.Context{
		MainModule: "",
		MainPkg:    pathToMainPkg,
		Modules: map[string]compiler.RawModule{
			"": rawMod, // TODO name?
		},
	})
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
	adapter proto.Adapter,
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
