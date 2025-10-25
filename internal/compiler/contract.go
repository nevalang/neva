package compiler

import (
	"context"

	"github.com/nevalang/neva/internal/compiler/ir"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

const (
	ExternDirective    src.Directive = "extern"
	BindDirective      src.Directive = "bind"
	AutoportsDirective src.Directive = "autoports"
)

type (
	Builder interface {
		Build(ctx context.Context, wd string) (RawBuild, string, *Error)
	}

	RawBuild struct {
		EntryModRef core.ModuleRef
		Modules     map[core.ModuleRef]RawModule
	}

	Parser interface {
		ParseModules(rawMods map[core.ModuleRef]RawModule) (map[core.ModuleRef]src.Module, *Error)
	}

	RawModule struct {
		Manifest src.ModuleManifest    // Manifest must be parsed by builder before passing into compiler
		Packages map[string]RawPackage // Packages themselves on the other hand can be parsed by compiler
	}

	RawPackage map[string][]byte

	Analyzer interface {
		AnalyzeBuild(mod src.Build, mainPkgName string) (src.Build, *Error)
	}

	Desugarer interface {
		Desugar(build src.Build) (src.Build, error)
	}

	Irgen interface {
		Generate(build src.Build, mainpkg string) (*ir.Program, error)
	}

	Backend interface {
		Emit(dst string, prog *ir.Program, trace bool) error
	}
)
