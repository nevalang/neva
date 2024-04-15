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

// Compile compiles given rawBuild to target language
// and uses specified backend to emit files to the destination.
func (c Compiler) Compile(
	src string,
	mainPkgName string,
	dstPath string,
) error {
	ir, err := c.CompileToIR(src, mainPkgName)
	if err != nil {
		return err
	}
	return c.backend.Emit(dstPath, ir)
}

func (c Compiler) CompileToIR(whereCLIExecuted string, whereEntryPkg string) (*ir.Program, *Error) {
	rawBuild, err := c.builder.Build(
		context.Background(),
		whereCLIExecuted,
		whereEntryPkg,
	)
	if err != nil {
		return nil, Error{
			Location: &sourcecode.Location{
				PkgName: whereEntryPkg,
			},
		}.Wrap(err)
	}

	parsedMods, err := c.parser.ParseModules(rawBuild.Modules)
	if err != nil {
		return nil, err
	}

	parsedBuild := sourcecode.Build{
		EntryModRef: rawBuild.EntryModRef,
		Modules:     parsedMods,
	}

	whereEntryPkg = strings.TrimPrefix(whereEntryPkg, "./")

	analyzedBuild, err := c.analyzer.AnalyzeExecutableBuild(parsedBuild, whereEntryPkg)
	if err != nil {
		return nil, err
	}

	desugaredBuild, err := c.desugarer.Desugar(analyzedBuild)
	if err != nil {
		return nil, err
	}

	irProg, err := c.irgen.Generate(desugaredBuild, whereEntryPkg)
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
