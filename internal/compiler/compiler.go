package compiler

import (
	"context"
	"fmt"
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

	if strings.HasPrefix(mainPkgName, "./") {
		mainPkgName = strings.TrimPrefix(mainPkgName, "./")
	}

	analyzedBuild, err := c.analyzer.AnalyzeExecutableBuild(parsedBuild, mainPkgName)
	if err != nil {
		return nil, err
	}

	desugaredBuild, err := c.desugarer.Desugar(analyzedBuild)
	if err != nil {
		return nil, err
	}

	irProg, err := c.irgen.Generate(ctx, desugaredBuild, mainPkgName)
	if err != nil {
		return nil, err
	}

	// FIXME no structBuilder func call in resulting irprog

	fmt.Println(
		JSONDump(parsedBuild.Modules[desugaredBuild.EntryModRef].Packages["struct_builder/with_sugar"]),
	)
	fmt.Println()
	fmt.Println()
	fmt.Println(
		JSONDump(analyzedBuild.Modules[desugaredBuild.EntryModRef].Packages["struct_builder/with_sugar"]),
	)
	fmt.Println()
	fmt.Println()
	fmt.Println(
		JSONDump(desugaredBuild.Modules[desugaredBuild.EntryModRef].Packages["struct_builder/with_sugar"]),
	)
	fmt.Println()
	fmt.Println()
	fmt.Println(
		JSONDump(irProg),
	)

	return irProg, nil
}

// New creates new Compiler instance.
// You can omit irgen and backend if all you need is Analyze method.
func New(
	parser Parser,
	desugarer Desugarer,
	analyzer Analyzer,
	irgen IRGenerator,
	backend Backend,
) Compiler {
	return Compiler{
		parser:    parser,
		desugarer: desugarer,
		analyzer:  analyzer,
		irgen:     irgen,
		backend:   backend,
	}
}
