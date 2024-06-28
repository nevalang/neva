package compiler

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/runtime/ir"
)

type Compiler struct {
	builder   Builder
	parser    Parser
	desugarer Desugarer
	analyzer  Analyzer
	irgen     IRGenerator
	backend   Backend
}

func (c Compiler) Compile(main string, dst string) error {
	ir, err := c.CompileToIR(main)
	if err != nil {
		return err
	}
	return c.backend.Emit(dst, ir)
}

func (c Compiler) CompileToIR(main string) (*ir.Program, *Error) {
	raw, root, err := c.builder.Build(context.Background(), main)
	if err != nil {
		return nil, Error{Location: &sourcecode.Location{PkgName: main}}.Wrap(err)
	}

	parsedMods, err := c.parser.ParseModules(raw.Modules)
	if err != nil {
		return nil, err
	}

	parsedBuild := sourcecode.Build{
		EntryModRef: raw.EntryModRef,
		Modules:     parsedMods,
	}

	main = strings.TrimPrefix(main, "./")
	main = strings.TrimPrefix(main, root+"/")

	analyzedBuild, err := c.analyzer.AnalyzeExecutableBuild(
		parsedBuild,
		main,
	)
	if err != nil {
		return nil, err
	}

	desugaredBuild, err := c.desugarer.Desugar(analyzedBuild)
	if err != nil {
		return nil, err
	}

	irProg, err := c.irgen.Generate(desugaredBuild, main, false)
	if err != nil {
		return nil, err
	}

	return irProg, nil
}

func New(
	builder Builder,
	parser Parser,
	desugarer Desugarer,
	analyzer Analyzer,
	irgen IRGenerator,
	backend Backend,
) Compiler {
	return Compiler{
		builder:   builder,
		parser:    parser,
		desugarer: desugarer,
		analyzer:  analyzer,
		irgen:     irgen,
		backend:   backend,
	}
}
