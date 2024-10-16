package compiler

import (
	"context"

	"github.com/nevalang/neva/internal/compiler/ir"
	src "github.com/nevalang/neva/internal/compiler/sourcecode"
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
		EntryModRef src.ModuleRef
		Modules     map[src.ModuleRef]RawModule
	}

	Parser interface {
		ParseModules(rawMods map[src.ModuleRef]RawModule) (map[src.ModuleRef]src.Module, *Error)
	}

	RawModule struct {
		Manifest src.ModuleManifest    // Manifest must be parsed by builder before passing into compiler
		Packages map[string]RawPackage // Packages themselves on the other hand can be parsed by compiler
	}

	RawPackage map[string][]byte

	Analyzer interface {
		AnalyzeExecutableBuild(mod src.Build, mainPkgName string) (src.Build, *Error)
	}

	Desugarer interface {
		Desugar(build src.Build) (src.Build, *Error)
	}

	IRGenerator interface {
		Generate(build src.Build, mainpkg string) (*ir.Program, *Error)
	}

	Backend interface {
		Emit(dst string, prog *ir.Program, trace bool) error
	}
)
