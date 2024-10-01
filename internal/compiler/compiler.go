package compiler

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/nevalang/neva/internal/compiler/sourcecode"
)

type Compiler struct {
	builder   Builder
	parser    Parser
	desugarer Desugarer
	analyzer  Analyzer
	irgen     IRGenerator
	backend   Backend
}

func (c Compiler) Compile(main string, output string, trace bool) error {
	result, err := c.CompileToIR(main, trace)
	if err != nil {
		return err
	}
	return c.backend.Emit(output, result.IR, trace)
}

type CompileResult struct {
	ParsedBuild    sourcecode.Build
	AnalyzedBuild  sourcecode.Build
	DesugaredBuild sourcecode.Build
	IR             *ir.Program
}

func (c Compiler) CompileToIR(main string, trace bool) (CompileResult, *Error) {
	raw, root, err := c.builder.Build(context.Background(), main)
	if err != nil {
		return CompileResult{}, Error{Location: &sourcecode.Location{PkgName: main}}.Wrap(err)
	}

	parsedMods, err := c.parser.ParseModules(raw.Modules)
	if err != nil {
		return CompileResult{}, err
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
		return CompileResult{}, err
	}

	desugaredBuild, err := c.desugarer.Desugar(analyzedBuild)
	if err != nil {
		return CompileResult{}, err
	}

	irProg, err := c.irgen.Generate(desugaredBuild, main, !trace)
	if err != nil {
		return CompileResult{}, err
	}

	return CompileResult{
		ParsedBuild:    parsedBuild,
		AnalyzedBuild:  analyzedBuild,
		DesugaredBuild: desugaredBuild,
		IR:             irProg,
	}, nil
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
