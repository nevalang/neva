package interpreter

import (
	"context"

	"github.com/nevalang/neva/internal/compiler/std"
	"github.com/nevalang/neva/internal/runtime"
	"github.com/nevalang/neva/internal/src"
	"github.com/nevalang/neva/pkg/ir"
)

type Interpreter struct {
	parser  SourceCodeParser
	irgen   IRGenerator
	rtgen   RuntimeProgramGenerator
	runtime Runtime
}

type SourceCodeParser interface {
	ParseFiles(context.Context, map[string][]byte) (map[string]src.Package, error)
	ParseFile(context.Context, []byte) (src.Package, error)
}

type IRGenerator interface {
	Generate(context.Context, map[string]src.Package) (*ir.Program, error)
}

type RuntimeProgramGenerator interface {
	Transform(context.Context, *ir.Program) (runtime.Program, error)
}

type Runtime interface {
	Run(context.Context, runtime.Program) (code int, err error)
}

func (i Interpreter) Interpret(ctx context.Context, bb []byte) (int, error) {
	singleFilePkg, err := i.parser.ParseFile(ctx, bb)
	if err != nil {
		return 0, err
	}

	prog := map[string]src.Package{
		"main": singleFilePkg,
		"std":  std.New(),
	}

	ll, err := i.irgen.Generate(ctx, prog)
	if err != nil {
		return 0, err
	}

	rprog, err := i.rtgen.Transform(ctx, ll)
	if err != nil {
		return 0, err
	}

	code, err := i.runtime.Run(ctx, rprog)
	if err != nil {
		return 0, err
	}

	return code, nil
}

func MustNew(
	parser SourceCodeParser,
	irgen IRGenerator,
	transformer RuntimeProgramGenerator,
	runtime Runtime,
) Interpreter {
	if parser == nil || irgen == nil || transformer == nil || runtime == nil {
		panic("nil argument")
	}
	return Interpreter{
		parser:  parser,
		irgen:   irgen,
		rtgen:   transformer,
		runtime: runtime,
	}
}
