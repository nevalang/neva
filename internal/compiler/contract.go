package compiler

import (
	"context"

	"github.com/nevalang/neva/internal/compiler/ir"
	src "github.com/nevalang/neva/internal/compiler/ast"
	"github.com/nevalang/neva/internal/compiler/ast/core"
)

const (
	ExternDirective    src.Directive = "extern"
	BindDirective      src.Directive = "bind"
	AutoportsDirective src.Directive = "autoports"
)

type (
	Builder interface {
		Build(ctx context.Context, workdir string) (RawBuild, string, *Error)
	}
	RawBuild struct {
		EntryModRef core.ModuleRef
		Modules     map[core.ModuleRef]RawModule
	}
)

type (
	Parser interface {
		ParseModules(rawMods map[core.ModuleRef]RawModule) (map[core.ModuleRef]src.Module, *Error)
	}
	RawModule struct {
		Manifest src.ModuleManifest    // Manifest must be parsed by builder before passing into compiler
		Packages map[string]RawPackage // Packages themselves on the other hand can be parsed by compiler
	}
	RawPackage map[string][]byte
)

type Analyzer interface {
	Analyze(mod src.Build, mainPkgName string) (src.Build, *Error)
}

type Desugarer interface {
	Desugar(build src.Build) (src.Build, error)
}

type Irgen interface {
	Generate(build src.Build, mainPkgName string) (*ir.Program, error)
}

type Backend interface {
	Emit(dst string, prog *ir.Program, isTraceEnabled bool) error
}
