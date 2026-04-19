//nolint:all // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
package compiler

import (
	"context"

	"github.com/nevalang/neva/internal/compiler/ir"
	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
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
		Modules     map[core.ModuleRef]RawModule
		EntryModRef core.ModuleRef
	}
)

type (
	Parser interface {
		ParseModules(rawMods map[core.ModuleRef]RawModule) (map[core.ModuleRef]src.Module, *Error)
	}
	RawModule struct {
		Packages map[string]RawPackage
		Manifest src.ModuleManifest
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
	GenerateForComponent(build src.Build, pkgName, componentName string) (*ir.Program, error)
}

type Backend interface {
	EmitExecutable(dst string, prog *ir.Program, trace bool) error
	EmitLibrary(dst string, exports []LibraryExport, trace bool) error
}

type LibraryExport struct {
	Program   *ir.Program
	Name      string
	Component src.Component
}
