package compiler

import (
	"context"
	"strings"

	"github.com/nevalang/neva/pkg/ir"
	src "github.com/nevalang/neva/pkg/sourcecode"
)

type Compiler struct {
	parser    Parser
	desugarer Desugarer
	analyzer  Analyzer
	irgen     IRGenerator
	backend   Backend
}

func (c Compiler) Compile(
	ctx context.Context,
	rawBuild RawBuild,
	workdirPath string,
	mainPkgName string,
) ([]byte, error) {
	ir, err := c.CompileToIR(ctx, rawBuild, workdirPath, mainPkgName)
	if err != nil {
		return nil, err
	}
	return c.backend.GenerateTarget(ir)
}

func (c Compiler) CompileToIR(
	ctx context.Context,
	rawBuild RawBuild,
	workdirPath string,
	mainPkgName string,
) (*ir.Program, *Error) {
	parsedMods, err := c.parser.ParseModules(rawBuild.Modules)
	if err != nil {
		return nil, err
	}

	parsedBuild := src.Build{
		EntryModRef: rawBuild.EntryModRef,
		Modules:     parsedMods,
	}

	desugaredBuild, err := c.desugarer.Desugar(parsedBuild)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(mainPkgName, "./") {
		mainPkgName = strings.TrimPrefix(mainPkgName, "./")
	}

	analyzedBuild, err := c.analyzer.AnalyzeExecutableBuild(desugaredBuild, mainPkgName)
	if err != nil {
		return nil, err
	}

	irProg, err := c.irgen.Generate(ctx, analyzedBuild, mainPkgName)
	if err != nil {
		return nil, err
	}

	return irProg, nil
}

// New creates new Compiler instance. You can omit irgen if all you need is Analyze method.
func New(
	parser Parser,
	desugarer Desugarer,
	analyzer Analyzer,
	irgen IRGenerator,
) Compiler {
	return Compiler{
		parser:    parser,
		desugarer: desugarer,
		analyzer:  analyzer,
		irgen:     irgen,
	}
}
