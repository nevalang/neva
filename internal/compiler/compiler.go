package compiler

import (
	"context"
	"fmt"
	"strings"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/pkg/ir"
)

type Compiler struct {
	parser   Parser
	analyzer Analyzer
	irgen    IRGen
}

type (
	Build struct {
		EntryModule string
		Modules     map[string]RawModule
	}

	RawModule struct {
		Manifest src.Manifest          // Manifest must be parsed by builder before passing into compiler
		Packages map[string]RawPackage // Packages themselves on the other hand can be parsed by compiler
	}

	Parser interface {
		ParsePackages(context.Context, map[string]RawPackage) (map[string]src.Package, error)
	}

	RawPackage map[string][]byte

	Analyzer interface {
		AnalyzeExecutable(prog src.Module, mainPkg string) (src.Module, error)
	}

	IRGen interface {
		Generate(ctx context.Context, mod src.Module, mainPkgName string) (*ir.Program, error)
	}
)

func (c Compiler) Compile(
	ctx context.Context,
	build Build,
	workdirPath string,
	mainPkgName string,
) (*ir.Program, error) {
	rawMod := build.Modules[build.EntryModule] // TODO support multimodule compilation

	if strings.HasPrefix(mainPkgName, "./") {
		mainPkgName = strings.TrimPrefix(mainPkgName, "./")
	}

	parsedPackages, err := c.parser.ParsePackages(ctx, rawMod.Packages)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	mod := src.Module{
		Manifest: rawMod.Manifest,
		Packages: parsedPackages,
	}

	analyzedProg, err := c.analyzer.AnalyzeExecutable(mod, mainPkgName)
	if err != nil {
		return nil, fmt.Errorf("analyzer: %w", err)
	}

	irProg, err := c.irgen.Generate(ctx, analyzedProg, mainPkgName)
	if err != nil {
		return nil, fmt.Errorf("generate IR: %w", err)
	}

	return irProg, nil
}

// New creates new Compiler instance. You can omit irgen if all you need is Analyze method.
func New(
	parser Parser,
	analyzer Analyzer,
	irgen IRGen,
) Compiler {
	return Compiler{
		parser:   parser,
		analyzer: analyzer,
		irgen:    irgen,
	}
}
