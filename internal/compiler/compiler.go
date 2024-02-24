package compiler

import (
	"context"
	"strings"

	"github.com/nevalang/neva/internal/runtime/ir"
	"github.com/nevalang/neva/pkg/sourcecode"
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

// CompileToIR compiles to intermediate representation
func (c Compiler) CompileToIR(
	src string,
	mainPkgName string,
) (*ir.Program, *Error) {
	rawBuild, err := c.builder.Build(context.Background(), src)
	if err != nil {
		return nil, Error{
			Location: &sourcecode.Location{
				PkgName: mainPkgName,
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

	mainPkgName = strings.TrimPrefix(mainPkgName, "./")

	analyzedBuild, err := c.analyzer.AnalyzeExecutableBuild(parsedBuild, mainPkgName)
	if err != nil {
		return nil, err
	}

	desugaredBuild, err := c.desugarer.Desugar(analyzedBuild)
	if err != nil {
		return nil, err
	}

	irProg, err := c.irgen.Generate(desugaredBuild, mainPkgName)
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
