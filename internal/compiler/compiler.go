package compiler

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/nevalang/neva/internal/compiler/sourcecode"
)

type Compiler struct {
	frontend  Frontend
	middleend Middleend
	backend   Backend
}

func (c Compiler) Compile(ctx context.Context, main string, output string, trace bool) error {
	feResult, err := c.frontend.Process(ctx, main)
	if err != nil {
		return err
	}

	meResult, err := c.middleend.Process(feResult)
	if err != nil {
		return err
	}

	return c.backend.Emit(output, meResult.IR, trace)
}

type Frontend struct {
	builder Builder
	parser  Parser
}

type FrontendResult struct {
	Root        string
	RawBuild    RawBuild
	ParsedBuild sourcecode.Build
}

func (f Frontend) Process(ctx context.Context, main string) (FrontendResult, *Error) {
	raw, root, err := f.builder.Build(ctx, main)
	if err != nil {
		return FrontendResult{}, Error{Location: &sourcecode.Location{PkgName: main}}.Wrap(err)
	}

	parsedMods, err := f.parser.ParseModules(raw.Modules)
	if err != nil {
		return FrontendResult{}, err
	}

	parsedBuild := sourcecode.Build{
		EntryModRef: raw.EntryModRef,
		Modules:     parsedMods,
	}

	return FrontendResult{
		ParsedBuild: parsedBuild,
		RawBuild:    raw,
		Root:        root,
	}, nil
}

type Middleend struct {
	desugarer Desugarer
	analyzer  Analyzer
	irgen     IRGenerator
}

type MiddleendResult struct {
	AnalyzedBuild  sourcecode.Build
	DesugaredBuild sourcecode.Build
	IR             *ir.Program
}

func (m Middleend) Process(feResult FrontendResult) (MiddleendResult, *Error) {
	main := strings.TrimPrefix(feResult.Root, "./")
	main = strings.TrimPrefix(main, feResult.Root+"/")

	analyzedBuild, err := m.analyzer.AnalyzeExecutableBuild(
		feResult.ParsedBuild,
		main,
	)
	if err != nil {
		return MiddleendResult{}, err
	}

	desugaredBuild, err := m.desugarer.Desugar(analyzedBuild)
	if err != nil {
		return MiddleendResult{}, err
	}

	irProg, err := m.irgen.Generate(desugaredBuild, main)
	if err != nil {
		return MiddleendResult{}, err
	}

	return MiddleendResult{
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
		frontend: Frontend{
			builder: builder,
			parser:  parser,
		},
		middleend: Middleend{
			desugarer: desugarer,
			analyzer:  analyzer,
			irgen:     irgen,
		},
		backend: backend,
	}
}
