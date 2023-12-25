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

type (
	Parser interface {
		ParseModules(rawMods map[src.ModuleRef]RawModule) (map[src.ModuleRef]src.Module, error)
	}

	RawPackage map[string][]byte

	Desugarer interface {
		Desugar(build src.Build) (src.Build, error)
	}

	Analyzer interface {
		AnalyzeExecutableBuild(mod src.Build, mainPkgName string) (src.Build, error)
	}

	IRGenerator interface {
		Generate(ctx context.Context, build src.Build, mainPkgName string) (*ir.Program, error)
	}

	RawBuild struct {
		EntryModRef src.ModuleRef
		Modules     map[src.ModuleRef]RawModule
	}

	RawModule struct {
		Manifest src.ModuleManifest    // Manifest must be parsed by builder before passing into compiler
		Packages map[string]RawPackage // Packages themselves on the other hand can be parsed by compiler
	}

	Backend interface {
		GenerateTarget(*ir.Program) ([]byte, error)
	}
)

// Compiler directives that dependency interface implementations must support.
const (
	RuntimeFuncDirective    src.Directive = "runtime_func"
	RuntimeFuncMsgDirective src.Directive = "runtime_func_msg"
)

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
) (*ir.Program, error) {
	parsedMods, err := c.parser.ParseModules(rawBuild.Modules)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	parsedBuild := src.Build{
		EntryModRef: rawBuild.EntryModRef,
		Modules:     parsedMods,
	}

	desugaredBuild, err := c.desugarer.Desugar(parsedBuild)
	if err != nil {
		return nil, fmt.Errorf("analyzer: %w", err)
	}

	if strings.HasPrefix(mainPkgName, "./") {
		mainPkgName = strings.TrimPrefix(mainPkgName, "./")
	}

	analyzedBuild, err := c.analyzer.AnalyzeExecutableBuild(desugaredBuild, mainPkgName)
	if err != nil {
		return nil, fmt.Errorf("analyzer: %w", err)
	}

	irProg, err := c.irgen.Generate(ctx, analyzedBuild, mainPkgName)
	if err != nil {
		return nil, fmt.Errorf("generate IR: %w", err)
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
